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

// Code generated by "stringer -type=ExecutorType"; DO NOT EDIT.

package parser

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnsupportExecutor-1]
	_ = x[InsertExecutor-2]
	_ = x[UpdateExecutor-3]
	_ = x[DeleteExecutor-4]
	_ = x[ReplaceIntoExecutor-5]
	_ = x[InsertOnDuplicateExecutor-6]
}

const _ExecutorType_name = "UnsupportExecutorInsertExecutorUpdateExecutorDeleteExecutorReplaceIntoExecutorInsertOnDuplicateExecutor"

var _ExecutorType_index = [...]uint8{0, 17, 31, 45, 59, 78, 103}

func (i ExecutorType) String() string {
	i -= 1
	if i < 0 || i >= ExecutorType(len(_ExecutorType_index)-1) {
		return "ExecutorType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ExecutorType_name[_ExecutorType_index[i]:_ExecutorType_index[i+1]]
}
