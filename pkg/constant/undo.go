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

package constant

const (
	DeleteFrom                     = "DELETE FROM "
	DefaultTransactionUndoLogTable = " undo_log "
	// UndoLogTableName Todo get from config
	UndoLogTableName  = DefaultTransactionUndoLogTable
	DeleteUndoLogSql  = DeleteFrom + UndoLogTableName + " WHERE " + UndoLogBranchXid + " = ? AND " + UndoLogXid + " = ?"
	SelectUndoLogSql  = "SELECT `log_status`,`context`,`rollback_info` FROM " + UndoLogTableName + " WHERE " + UndoLogBranchXid + " = ? AND " + UndoLogXid + " = ? FOR UPDATE"
	DeleteSqlTemplate = "DELETE FROM %s WHERE %s "
	// InsertSqlTemplate INSERT INTO a (x, y, z, pk) VALUES (?, ?, ?, ?)
	InsertSqlTemplate = "INSERT INTO %s (%s) VALUES (%s)"
	// UpdateSqlTemplate UPDATE a SET x = ?, y = ?, z = ? WHERE pk1 in (?) pk2 in (?)
	UpdateSqlTemplate = "UPDATE %s SET %s WHERE %s "
)

// undo log status
const (
	// UndoLogStatusNormal This state can be properly rolled back by services
	UndoLogStatusNormal = iota
	// UndoLogStatusGlobalFinished This state prevents the branch transaction from inserting undo_log after the global transaction is rolled back.
	UndoLogStatusGlobalFinished
)

// undo log compress
const (
	CompressorTypeKey = "compressorTypeKey"
	SerializerKey     = "serializerKey"
)

// table schema
const (
	IndexSchemaSql = "SELECT `INDEX_NAME`, `COLUMN_NAME`, `NON_UNIQUE`, `INDEX_TYPE`, `COLLATION`, `CARDINALITY` " +
		"FROM `INFORMATION_SCHEMA`.`STATISTICS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"

	ColumnSchemaSql = "select TABLE_CATALOG, TABLE_NAME, TABLE_SCHEMA, COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, COLUMN_KEY, " +
		" IS_NULLABLE, EXTRA from INFORMATION_SCHEMA.COLUMNS where `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"
)

const DBName = "seata"
