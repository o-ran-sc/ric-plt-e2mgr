package dispatcher

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"sync"
	"time"
	"xappmock/enums"
	"xappmock/models"
	"xappmock/rmr"
	"xappmock/sender"
)

// Id -> Command
var configuration = make(map[string]*models.JsonCommand)
// Rmr Message Id -> Command
var waitForRmrMessageType = make(map[int]*models.JsonCommand)

func addRmrMessageToWaitFor(rmrMessageToWaitFor string, command models.JsonCommand) error {
	rmrMsgId, err := rmr.MessageIdToUint(rmrMessageToWaitFor)

	if err != nil {
		return errors.New(fmt.Sprintf("invalid rmr message id: %s", rmrMessageToWaitFor))
	}

	waitForRmrMessageType[int(rmrMsgId)] = &command
	return nil
}

type Dispatcher struct {
	rmrService    *rmr.Service
	processResult models.ProcessResult
}

func (d *Dispatcher) GetProcessResult() models.ProcessResult {
	return d.processResult
}

func New(rmrService *rmr.Service) *Dispatcher {
	return &Dispatcher{
		rmrService: rmrService,
	}
}

func (d *Dispatcher) JsonCommandsDecoderCB(cmd models.JsonCommand) error {
	if len(cmd.Id) == 0 {
		return errors.New(fmt.Sprintf("invalid cmd, no id"))
	}
	configuration[cmd.Id] = &cmd

	if len(cmd.ReceiveRmrMessageType) == 0 {
		return nil
	}

	return addRmrMessageToWaitFor(cmd.ReceiveRmrMessageType, cmd)
}

func (d *Dispatcher) sendNoRepeat(command models.JsonCommand) {
	err := sender.SendJsonRmrMessage(command, nil, d.rmrService)

	if err != nil {
		log.Printf("Dispatcher.sendHandler - error sending rmr message: %s", err)
		d.processResult.Err = err
		d.processResult.Stats.SentErrorCount++
		return
	}

	d.processResult.Stats.SentCount++

}

func (d *Dispatcher) sendWithRepeat(ctx context.Context, command models.JsonCommand) {
	for repeatCount := command.RepeatCount; repeatCount > 0; repeatCount-- {

		select {
		case <-ctx.Done():
			return
		default:
		}

		err := sender.SendJsonRmrMessage(command, nil, d.rmrService)

		if err != nil {
			log.Printf("Dispatcher.sendHandler - error sending rmr message: %s", err)
			d.processResult.Stats.SentErrorCount++
			continue
		}

		d.processResult.Stats.SentCount++
		time.Sleep(time.Duration(command.RepeatDelayInMs) * time.Millisecond)
	}
}

func (d *Dispatcher) sendHandler(ctx context.Context, sendAndReceiveWg *sync.WaitGroup, command models.JsonCommand) {

	defer sendAndReceiveWg.Done()
	var listenAndHandleWg sync.WaitGroup

	if len(command.ReceiveRmrMessageType) > 0 {
		err := addRmrMessageToWaitFor(command.ReceiveRmrMessageType, command)

		if err != nil {
			d.processResult.Err = err
			return
		}

		listenAndHandleWg.Add(1)
		go d.listenAndHandle(ctx, &listenAndHandleWg, command.RepeatCount)
	}

	if command.RepeatCount == 0 {
		d.sendNoRepeat(command)
	} else {
		d.sendWithRepeat(ctx, command)
	}

	if len(command.ReceiveRmrMessageType) > 0 {
		listenAndHandleWg.Wait()
	}
}

func (d *Dispatcher) receiveHandler(ctx context.Context, sendAndReceiveWg *sync.WaitGroup, command models.JsonCommand) {

	defer sendAndReceiveWg.Done()

	err := addRmrMessageToWaitFor(command.ReceiveRmrMessageType, command)

	if err != nil {
		d.processResult.Err = err
		return
	}

	var listenAndHandleWg sync.WaitGroup
	listenAndHandleWg.Add(1) // this is due to the usage of listenAndHandle as a goroutine in the sender case
	d.listenAndHandle(ctx, &listenAndHandleWg, command.RepeatCount)
}

func getMergedCommand(cmd *models.JsonCommand) (models.JsonCommand, error) {
	var command models.JsonCommand
	if len(cmd.Id) == 0 {
		return command, errors.New(fmt.Sprintf("invalid command, no id"))
	}

	command = *cmd

	conf, ok := configuration[cmd.Id]

	if ok {
		command = *conf
		mergeConfigurationAndCommand(&command, cmd)
	}

	return command, nil
}

func (d *Dispatcher) ProcessJsonCommand(ctx context.Context, cmd *models.JsonCommand) {

	command, err := getMergedCommand(cmd)

	if err != nil {
		d.processResult.Err = err
		return
	}

	var sendAndReceiveWg sync.WaitGroup

	commandAction := enums.CommandAction(command.Action)

	switch commandAction {

	case enums.SendRmrMessage:
		sendAndReceiveWg.Add(1)
		go d.sendHandler(ctx, &sendAndReceiveWg, command)
	case enums.ReceiveRmrMessage:
		sendAndReceiveWg.Add(1)
		go d.receiveHandler(ctx, &sendAndReceiveWg, command)
	default:
		d.processResult = models.ProcessResult{Err: errors.New(fmt.Sprintf("invalid command action %s", command.Action))}
		return
	}

	sendAndReceiveWg.Wait()
}

func (d *Dispatcher) listenAndHandleNoRepeat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mbuf, err := d.rmrService.RecvMessage()

		if err != nil {
			d.processResult.Err = err
			d.processResult.Stats.ReceivedErrorCount++
			return
		}

		_, ok := waitForRmrMessageType[mbuf.MType]

		if !ok {
			log.Printf("#Dispatcher.listenAndHandle - Unexpected msg: %s", mbuf)
			d.processResult.Stats.ReceivedUnexpectedCount++
			continue
		}

		log.Printf("#Dispatcher.listenAndHandle - expected msg: %s", mbuf)
		d.processResult.Stats.ReceivedExpectedCount++
		return
	}
}

func (d *Dispatcher) receive(ctx context.Context) {

}

func (d *Dispatcher) listenAndHandleWithRepeat(ctx context.Context, repeatCount int) {
	for d.processResult.Stats.ReceivedExpectedCount < repeatCount {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mbuf, err := d.rmrService.RecvMessage()

		if err != nil {
			log.Printf("#Dispatcher.listenAndHandle - error receiving message: %s", err)
			d.processResult.Stats.ReceivedErrorCount++
			continue
		}

		_, ok := waitForRmrMessageType[mbuf.MType]

		if !ok {
			log.Printf("#Dispatcher.listenAndHandle - Unexpected msg: %s", mbuf)
			d.processResult.Stats.ReceivedUnexpectedCount++
			continue
		}

		log.Printf("#Dispatcher.listenAndHandle - expected msg: %s", mbuf)
		d.processResult.Stats.ReceivedExpectedCount++
	}
}

func (d *Dispatcher) listenAndHandle(ctx context.Context, listenAndHandleWg *sync.WaitGroup, repeatCount int) {

	defer listenAndHandleWg.Done()

	if repeatCount == 0 {
		d.listenAndHandleNoRepeat(ctx)
		return
	}

	d.listenAndHandleWithRepeat(ctx, repeatCount)
}

func mergeConfigurationAndCommand(conf *models.JsonCommand, cmd *models.JsonCommand) {
	nFields := reflect.Indirect(reflect.ValueOf(cmd)).NumField()

	for i := 0; i < nFields; i++ {
		if fieldValue := reflect.Indirect(reflect.ValueOf(cmd)).Field(i); fieldValue.IsValid() {
			switch fieldValue.Kind() {
			case reflect.String:
				if fieldValue.Len() > 0 {
					reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if fieldValue.Int() != 0 {
					reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
				}
			case reflect.Bool:
				if fieldValue.Bool() {
					reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
				}
			case reflect.Float64, reflect.Float32:
				if fieldValue.Float() != 0 {
					reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
				}
			default:
				reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
			}
		}
	}
}
