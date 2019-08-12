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

package sessions

import "time"
import "e2mgr/models"
/*
 * Container for session data
 * Note:
 * - If this is the only item in the session data, we should eliminate this session manager
 *   and just send the start time as part of the xaction field in the RMR message.
 */
type E2SessionDetails struct {
	SessionStart time.Time
	Request *models.RequestDetails
}

/*
 * Map of session id to session details.
 * Notes:
 * - Since the transaction id is the CellId, there is no
 *   need to delete the entry when a response is received nor
 *   provide a mechanism for removing stale entries (no response received).
 *   Having said that, deleting the entry on a successful flow may still be a good idea
 *   in order to avoid pinning large amount of memory (help the GC).
 *
 * TODO:
 *  - Synchronize access.
 */
type E2Sessions map[string]E2SessionDetails
