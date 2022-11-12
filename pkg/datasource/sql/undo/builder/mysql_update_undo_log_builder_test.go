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

package builder

import (
	"context"
	"database/sql/driver"
	"github.com/agiledragon/gomonkey"
	"github.com/seata/seata-go/pkg/datasource/sql/datasource/mysql"
	"github.com/seata/seata-go/pkg/datasource/sql/types"
	"reflect"
	"testing"

	"github.com/seata/seata-go/pkg/datasource/sql/parser"

	_ "github.com/arana-db/parser/test_driver"
	_ "github.com/seata/seata-go/pkg/util/log"
	"github.com/stretchr/testify/assert"
)

func TestBuildSelectSQLByUpdate(t *testing.T) {
	var (
		builder = MySQLUpdateUndoLogBuilder{}
	)
	//stub := gomonkey.ApplyMethod(reflect.TypeOf(mysql.GetTableMetaInstance()), "GetTableMeta", func(_ *datasource.TableMetaCache, ctx context.Context, dbName, tableName string, conn driver.Conn) (*types.TableMeta, error) {
	//	return &types.TableMeta{
	//
	//	}, nil
	//})
	stub := gomonkey.ApplyMethod(reflect.TypeOf(mysql.GetTableMetaInstance()), "GetTableMeta", func(_ *mysql.TableMetaCache, ctx context.Context, dbName, tableName string, conn driver.Conn) (*types.TableMeta, error) {
		return &types.TableMeta{
			Indexs: map[string]types.IndexMeta{
				"id": types.IndexMeta{
					ColumnName: "id",
					IType:      types.IndexTypePrimaryKey,
				},
			},
		}, nil
	})
	defer stub.Reset()

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
			expectQuery:     "SELECT SQL_NO_CACHE name,age,id FROM t_user WHERE id=? FOR UPDATE",
			expectQueryArgs: []driver.Value{100},
		},
		{
			sourceQuery:     "update t_user set name = ?, age = ? where id = ? and name = 'Jack' and age between ? and ?",
			sourceQueryArgs: []driver.Value{"Jack", 1, 100, 18, 28},
			expectQuery:     "SELECT SQL_NO_CACHE name,age,id FROM t_user WHERE id=? AND name=_UTF8MB4Jack AND age BETWEEN ? AND ? FOR UPDATE",
			expectQueryArgs: []driver.Value{100, 18, 28},
		},
		{
			sourceQuery:     "update t_user set name = ?, age = ? where id = ? and name = 'Jack' and age in (?,?)",
			sourceQueryArgs: []driver.Value{"Jack", 1, 100, 18, 28},
			expectQuery:     "SELECT SQL_NO_CACHE name,age,id FROM t_user WHERE id=? AND name=_UTF8MB4Jack AND age IN (?,?) FOR UPDATE",
			expectQueryArgs: []driver.Value{100, 18, 28},
		},
		{
			sourceQuery:     "update t_user set name = ?, age = ? where kk between ? and ? and id = ? and addr in(?,?) and age > ? order by name desc limit ?",
			sourceQueryArgs: []driver.Value{"Jack", 1, 10, 20, 17, "Beijing", "Guangzhou", 18, 2},
			expectQuery:     "SELECT SQL_NO_CACHE name,age,id FROM t_user WHERE kk BETWEEN ? AND ? AND id=? AND addr IN (?,?) AND age>? ORDER BY name DESC LIMIT ? FOR UPDATE",
			expectQueryArgs: []driver.Value{10, 20, 17, "Beijing", "Guangzhou", 18, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := parser.DoParser(tt.sourceQuery)
			assert.Nil(t, err)
			query, args, err := builder.buildBeforeImageSQL(context.Background(), &types.ExecContext{ParseContext: c}, tt.sourceQueryArgs)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectQuery, query)
			assert.Equal(t, tt.expectQueryArgs, args)
		})
	}
}
