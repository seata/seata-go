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

package at

import (
	"context"
	"database/sql/driver"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"seata.apache.org/seata-go/pkg/datasource/sql/exec"
	"seata.apache.org/seata-go/pkg/datasource/sql/parser"
	"seata.apache.org/seata-go/pkg/datasource/sql/types"
	"seata.apache.org/seata-go/pkg/datasource/sql/undo"
	"seata.apache.org/seata-go/pkg/datasource/sql/util"
	_ "seata.apache.org/seata-go/pkg/util/log"
)

var (
	MetaDataMap map[string]*types.TableMeta
)

func initTest() {
	MetaDataMap = map[string]*types.TableMeta{
		"t_user": {
			TableName: "t_user",
			Indexs: map[string]types.IndexMeta{
				"id": {
					IType: types.IndexTypePrimaryKey,
					Columns: []types.ColumnMeta{
						{ColumnName: "id"},
					},
				},
			},
			Columns: map[string]types.ColumnMeta{
				"id": {
					ColumnDef:  nil,
					ColumnName: "id",
				},
				"name": {
					ColumnDef:  nil,
					ColumnName: "name",
				},
				"age": {
					ColumnDef:  nil,
					ColumnName: "age",
				},
			},
			ColumnNames: []string{"id", "name", "age"},
		},
		"t1": {
			TableName: "t1",
			Indexs: map[string]types.IndexMeta{
				"id": {
					IType: types.IndexTypePrimaryKey,
					Columns: []types.ColumnMeta{
						{ColumnName: "id"},
					},
				},
			},
			Columns: map[string]types.ColumnMeta{
				"id": {
					ColumnDef:  nil,
					ColumnName: "id",
				},
				"name": {
					ColumnDef:  nil,
					ColumnName: "name",
				},
				"age": {
					ColumnDef:  nil,
					ColumnName: "age",
				},
			},
			ColumnNames: []string{"id", "name", "age"},
		},
		"t2": {
			TableName: "t2",
			Indexs: map[string]types.IndexMeta{
				"id": {
					IType: types.IndexTypePrimaryKey,
					Columns: []types.ColumnMeta{
						{ColumnName: "id"},
					},
				},
			},
			Columns: map[string]types.ColumnMeta{
				"id": {
					ColumnDef:  nil,
					ColumnName: "id",
				},
				"name": {
					ColumnDef:  nil,
					ColumnName: "name",
				},
				"age": {
					ColumnDef:  nil,
					ColumnName: "age",
				},
				"kk": {
					ColumnDef:  nil,
					ColumnName: "kk",
				},
				"addr": {
					ColumnDef:  nil,
					ColumnName: "addr",
				},
			},
			ColumnNames: []string{"id", "name", "age", "kk", "addr"},
		},
		"t3": {
			TableName: "t3",
			Indexs: map[string]types.IndexMeta{
				"id": {
					IType: types.IndexTypePrimaryKey,
					Columns: []types.ColumnMeta{
						{ColumnName: "id"},
					},
				},
			},
			Columns: map[string]types.ColumnMeta{
				"id": {
					ColumnDef:  nil,
					ColumnName: "id",
				},
				"age": {
					ColumnDef:  nil,
					ColumnName: "age",
				},
			},
			ColumnNames: []string{"id", "age"},
		},
		"t4": {
			TableName: "t4",
			Indexs: map[string]types.IndexMeta{
				"id": {
					IType: types.IndexTypePrimaryKey,
					Columns: []types.ColumnMeta{
						{ColumnName: "id"},
					},
				},
			},
			Columns: map[string]types.ColumnMeta{
				"id": {
					ColumnDef:  nil,
					ColumnName: "id",
				},
				"age": {
					ColumnDef:  nil,
					ColumnName: "age",
				},
			},
			ColumnNames: []string{"id", "age"},
		},
	}

	undo.InitUndoConfig(undo.Config{OnlyCareUpdateColumns: true})
}

func TestMain(m *testing.M) {
	// 调用初始化函数
	initTest()

	// 启动测试
	os.Exit(m.Run())
}

func TestBuildSelectSQLByUpdate(t *testing.T) {

	tests := []struct {
		name            string
		sourceQuery     string
		sourceQueryArgs []driver.Value
		expectQuery     string
		expectQueryArgs []driver.Value
	}{
		{
			sourceQuery:     "update t_user set name = ?, age = ? where id = ?",
			sourceQueryArgs: []driver.Value{"Jack", 1, 100},
			expectQuery:     "SELECT SQL_NO_CACHE name,age,t_user.id FROM t_user WHERE id=? FOR UPDATE",
			expectQueryArgs: []driver.Value{100},
		},
		{
			sourceQuery:     "update t_user set name = ?, age = ? where id = ? and name = 'Jack' and age between ? and ?",
			sourceQueryArgs: []driver.Value{"Jack", 1, 100, 18, 28},
			expectQuery:     "SELECT SQL_NO_CACHE name,age,t_user.id FROM t_user WHERE id=? AND name=_UTF8MB4Jack AND age BETWEEN ? AND ? FOR UPDATE",
			expectQueryArgs: []driver.Value{100, 18, 28},
		},
		{
			sourceQuery:     "update t_user set name = ?, age = ? where id = ? and name = 'Jack' and age in (?,?)",
			sourceQueryArgs: []driver.Value{"Jack", 1, 100, 18, 28},
			expectQuery:     "SELECT SQL_NO_CACHE name,age,t_user.id FROM t_user WHERE id=? AND name=_UTF8MB4Jack AND age IN (?,?) FOR UPDATE",
			expectQueryArgs: []driver.Value{100, 18, 28},
		},
		{
			sourceQuery:     "update t_user set name = ?, age = ? where kk between ? and ? and id = ? and addr in(?,?) and age > ? order by name desc limit ?",
			sourceQueryArgs: []driver.Value{"Jack", 1, 10, 20, 17, "Beijing", "Guangzhou", 18, 2},
			expectQuery:     "SELECT SQL_NO_CACHE name,age,t_user.id FROM t_user WHERE kk BETWEEN ? AND ? AND id=? AND addr IN (?,?) AND age>? ORDER BY name DESC LIMIT ? FOR UPDATE",
			expectQueryArgs: []driver.Value{10, 20, 17, "Beijing", "Guangzhou", 18, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := parser.DoParser(tt.sourceQuery)
			assert.Nil(t, err)
			executor := NewUpdateExecutor(c, &types.ExecContext{Values: tt.sourceQueryArgs, NamedValues: util.ValueToNamedValue(tt.sourceQueryArgs)}, []exec.SQLHook{})
			query, args, err := executor.(*updateExecutor).buildBeforeImageSQL(context.Background(), MetaDataMap["t_user"], util.ValueToNamedValue(tt.sourceQueryArgs))
			assert.Nil(t, err)
			assert.Equal(t, tt.expectQuery, query)
			assert.Equal(t, tt.expectQueryArgs, util.NamedValueToValue(args))
		})
	}
}
