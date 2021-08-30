package main

// /*
//  * Licensed to the Apache Software Foundation (ASF) under one
//  * or more contributor license agreements. See the NOTICE file
//  * distributed with this work for additional information
//  * regarding copyright ownership. The ASF licenses this file
//  * to you under the Apache License, Version 2.0 (the
//  * "License"); you may not use this file except in compliance
//  * with the License. You may obtain a copy of the License at
//  *
//  *   http://www.apache.org/licenses/LICENSE-2.0
//  *
//  * Unless required by applicable law or agreed to in writing,
//  * software distributed under the License is distributed on an
//  * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
//  * KIND, either express or implied. See the License for the
//  * specific language governing permissions and limitations
//  * under the License.
//  */

// import (
// 	"context"
// 	"thr/jmx"
// 	"thr/shared"
// )

// type CalculatorHandler struct {
// 	log map[int]*shared.SharedStruct
// }

// func NewCalculatorHandler() *CalculatorHandler {
// 	return &CalculatorHandler{log: make(map[int]*shared.SharedStruct)}
// }

// func (p *CalculatorHandler) Connect(ctx context.Context, config *jmx.JMXConfig) (r bool, err error) {
// 	// fmt.Print("ping()\n")
// 	return true, nil
// }

// func (p *CalculatorHandler) QueryMbean(ctx context.Context, beanName string) (r *jmx.Output, err error) {
// 	// fmt.Print("add(", num1, ",", num2, ")\n")
// 	return &jmx.Output{JSON: "{}"}, nil
// }

// func (p *CalculatorHandler) GetStruct(ctx context.Context, key int32) (*shared.SharedStruct, error) {
// 	// fmt.Print("getStruct(", key, ")\n")
// 	v, _ := p.log[int(key)]
// 	return v, nil
// }

// func (p *CalculatorHandler) Zip(ctx context.Context) (err error) {
// 	// fmt.Print("zip()\n")
// 	return nil
// }
