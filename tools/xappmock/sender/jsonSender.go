//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package sender

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode"
	"xappmock/models"
	"xappmock/rmr"
)

var counter uint64

func SendJsonRmrMessage(command models.JsonCommand /*the copy is modified locally*/, xAction *[]byte, r *rmr.Service) error {
	var payload []byte
	_, err := fmt.Sscanf(command.PackedPayload, "%x", &payload)
	if err != nil {
		return errors.New(fmt.Sprintf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err))
	}
	command.PackedPayload = string(payload)
	command.TransactionId = expandTransactionId(command.TransactionId)
	if len(command.TransactionId) == 0 {
		command.TransactionId = string(*xAction)
	}
	command.PayloadHeader = expandPayloadHeader(command.PayloadHeader, &command)
	log.Printf("#jsonSender.SendJsonRmrMessage - command payload header: %s", command.PayloadHeader)
	rmrMsgId, err := rmr.MessageIdToUint(command.RmrMessageType)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid rmr message id: %s", command.RmrMessageType))
	}

	msg := append([]byte(command.PayloadHeader), payload...)
	messageInfo := models.GetMessageInfoAsJson(int(rmrMsgId), command.Meid, msg, []byte(command.TransactionId))
	log.Printf("#rmr.Service.SendMessage - %s", messageInfo)

	_, err = r.SendMessage(int(rmrMsgId), command.Meid, msg, []byte(command.TransactionId))
	return err
}

/*
 * transactionId (xAction): The value may have a fixed value or $ or <prefix>$.
 * $ is replaced by a value generated at runtime (possibly unique per message sent).
 * If the tag does not exist, then the mock shall use the value taken from the incoming message.
 */
func expandTransactionId(id string) string {
	if len(id) == 1 && id[0] == '$' {
		return fmt.Sprintf("%d", incAndGetCounter())
	}
	if len(id) > 1 && id[len(id)-1] == '$' {
		return fmt.Sprintf("%s%d", id[:len(id)-1], incAndGetCounter())
	}
	return id
}

/*
 * payloadHeader: A prefix to combine with the payload that will be the message’s payload. The value may include variables of the format $<name> or #<name> where:
 *   $<name> expands to the value of <name> if it exists or the empty string if not.
 *   #<name> expands to the length of the value of <name> if it exists or omitted if not.
 * The intention is to allow the Mock to construct the payload header required by the setup messages (ranIp|ranPort|ranName|payload len|<payload>).
 * Example: “payloadHeader”: “$ranIp|$ranPort|$ranName|#packedPayload|”
 */

func expandPayloadHeader(header string, command *models.JsonCommand) string {
	var name strings.Builder
	var expandedHeader strings.Builder

	r := strings.NewReader(header)
	ch, err := r.ReadByte()
	for {
		if err != nil {
			break
		}

		switch ch {
		case '$':
			for {
				ch, err = r.ReadByte() //on error ch == 0
				if unicode.IsDigit(rune(ch)) || unicode.IsLetter(rune(ch)) {
					if name.Len() == 0 {
						name.WriteByte(byte(unicode.ToUpper(rune(ch))))
					} else {
						name.WriteByte(ch)
					}
				} else {
					if fieldValue := reflect.Indirect(reflect.ValueOf(command)).FieldByName(name.String()); fieldValue.IsValid() {
						switch fieldValue.Kind() {
						case reflect.String:
							expandedHeader.WriteString(fieldValue.String())
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							expandedHeader.WriteString(strconv.FormatInt(fieldValue.Int(), 10))
						case reflect.Bool:
							expandedHeader.WriteString(strconv.FormatBool(fieldValue.Bool()))
						case reflect.Float64, reflect.Float32:
							expandedHeader.WriteString(fmt.Sprintf("%g", fieldValue.Float()))
						default:
							log.Fatalf("#jsonSender.expandPayloadHeader - invalid type for $%s, value must be a string, an int, a bool or a float", name.String())
						}
					}
					name.Reset()
					break
				}
			}
		case '#':
			for {
				ch, err = r.ReadByte() //on error ch == 0
				if unicode.IsDigit(rune(ch)) || unicode.IsLetter(rune(ch)) {
					if name.Len() == 0 {
						name.WriteByte(byte(unicode.ToUpper(rune(ch))))
					} else {
						name.WriteByte(ch)
					}
				} else {
					if fieldValue := reflect.Indirect(reflect.ValueOf(command)).FieldByName(name.String()); fieldValue.IsValid() {
						if fieldValue.Kind() == reflect.String {
							expandedHeader.WriteString(strconv.FormatInt(int64(len(fieldValue.String())), 10))
						} else {
							log.Fatalf("#jsonSender.expandPayloadHeader - invalid type for #%s, value must be a string", name.String())
						}
					}
					name.Reset()
					break
				}
			}
		default:
			if unicode.IsPrint(rune(ch)) {
				expandedHeader.WriteByte(ch)
			}
			ch, err = r.ReadByte()
		}
	}
	return expandedHeader.String()
}

func incAndGetCounter() uint64 {
	return atomic.AddUint64(&counter, 1)
}

func init() {
	counter = uint64(time.Now().Unix() - 1572000000)
}
