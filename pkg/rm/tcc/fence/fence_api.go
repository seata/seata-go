/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zouyx/agollo/v3/component/log"

	"github.com/seata/seata-go/pkg/rm/tcc/fence/constant"
	"github.com/seata/seata-go/pkg/rm/tcc/fence/handler"
	"github.com/seata/seata-go/pkg/tm"
)

func WithFence(ctx context.Context, tx *sql.Tx, callback func() error) (resErr error) {
	fencePhase := tm.GetFencePhase(ctx)
	if fencePhase == constant.FencePhaseNotExist {
		return fmt.Errorf("xid %s, tx name %s, fence phase not exist", tm.GetXID(ctx), tm.GetTxName(ctx))
	}
	// deal with unknown situations
	defer func() {
		if err, ok := recover().(error); ok {
			if errRollback := tx.Rollback(); errRollback != nil {
				resErr = fmt.Errorf("rollback cause: %s, rollback err: %s", err.Error(), errRollback.Error())
			}
			log.Error(err)
			resErr = err
		}
	}()

	if fencePhase == constant.FencePhasePrepare {
		return handler.GetFenceHandlerSingleton().PrepareFence(ctx, tx, callback)
	}

	if fencePhase == constant.FencePhaseCommit {
		return handler.GetFenceHandlerSingleton().CommitFence(ctx, tx, callback)
	}

	if fencePhase == constant.FencePhaseRollback {
		return handler.GetFenceHandlerSingleton().RollbackFence(ctx, tx, callback)
	}
	return nil
}
