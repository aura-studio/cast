# 设计文档

## 概述

本设计文档描述了 `github.com/aura-studio/cast` 库的类型转换完整性分析和验证方案。目标是确保库支持的所有类型之间能够进行 N×N 的完整转换。

## 架构

### 当前架构

```
┌─────────────────────────────────────────────────────────────┐
│                      cast 库架构                             │
├─────────────────────────────────────────────────────────────┤
│  cast.go          - 核心辅助函数 (indirect, indirectToStringerOrError)  │
│  decimal.go       - 数值类型转换 (ToInt*, ToUint*, ToFloat*, ToBig*, ToComplex*, ToBool)  │
│  string.go        - 字符串类型转换 (ToString, ToBytes, ToStringer, ToError)  │
│  time.go          - 时间类型转换 (ToDuration, ToTimeZone)  │
└─────────────────────────────────────────────────────────────┘
```

### 类型转换矩阵

#### 源类型（25 种）
| 类别 | 类型 |
|------|------|
| 有符号整数 | int, int8, int16, int32, int64 |
| 无符号整数 | uint, uint8, uint16, uint32, uint64 |
| 浮点数 | float32, float64 |
| 大数 | *big.Int, *big.Float, *big.Rat |
| 复数 | complex64, complex128 |
| 布尔 | bool |
| 时间 | time.Duration, time.Location, *time.Location |
| 字符串 | string, []byte, fmt.Stringer, error |
| 空值 | nil |

#### 目标类型（22 种）
| 类别 | 类型 | 转换函数 |
|------|------|----------|
| 有符号整数 | int, int8, int16, int32, int64 | ToInt, ToInt8, ToInt16, ToInt32, ToInt64 |
| 无符号整数 | uint, uint8, uint16, uint32, uint64 | ToUint, ToUint8, ToUint16, ToUint32, ToUint64 |
| 浮点数 | float32, float64 | ToFloat32, ToFloat64 |
| 大数 | *big.Int, *big.Float, *big.Rat | ToBigInt, ToBigFloat, ToBigRat |
| 复数 | complex64, complex128 | ToComplex64, ToComplex128 |
| 布尔 | bool | ToBool |
| 时间 | time.Duration, *time.Location | ToDuration, ToTimeZone |
| 字符串 | string, []byte, fmt.Stringer, error | ToString, ToBytes, ToStringer, ToError |

## 组件和接口

### 转换函数接口模式

每个目标类型都有两个转换函数：

```go
// 简单版本 - 忽略错误，返回零值
func ToXxx(a any) Xxx {
    v, _ := ToXxxE(a)
    return v
}

// 带错误版本 - 返回结果和错误
func ToXxxE(a any) (Xxx, error) {
    a = indirectToStringerOrError(a)
    switch v := a.(type) {
    case Type1:
        // 转换逻辑
    case Type2:
        // 转换逻辑
    // ... 处理所有源类型
    case nil:
        return zeroValue, nil
    default:
        return zeroValue, fmt.Errorf("unable to cast %#v of type %T to Xxx", a, a)
    }
}
```

### 辅助函数

```go
// indirectToStringerOrError - 解引用指针，直到遇到 fmt.Stringer 或 error 接口
func indirectToStringerOrError(a any) any
```

## 数据模型

### 类型转换规则

| 源类型 → 目标类型 | 转换规则 |
|------------------|----------|
| 整数 → 整数 | 直接类型转换，可能截断 |
| 整数 → 无符号整数 | 负数返回错误 |
| 浮点 → 整数 | 截断小数部分 |
| 复数 → 实数 | 取实部 |
| 数值 → bool | 非零为 true |
| bool → 数值 | true=1, false=0 |
| 字符串 → 数值 | 解析字符串 |
| 数值 → 字符串 | 格式化为字符串 |
| Duration → 数值 | 纳秒数 |
| Location → 数值 | 时区偏移秒数 |

## 正确性属性

*正确性属性是一种应该在系统所有有效执行中保持为真的特征或行为——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: 整数类型相互转换一致性

*For any* 有符号整数值 v 和任意两个有符号整数类型 T1、T2，将 v 从 T1 转换到 T2 再转回 T1，如果 v 在 T2 的范围内，则结果应等于原值。

**Validates: Requirements 1.1, 1.2**

### Property 2: 负数转无符号整数返回错误

*For any* 负数值 v（有符号整数或浮点数），将其转换为任意无符号整数类型时，带 E 后缀的函数应返回非 nil 错误。

**Validates: Requirements 1.3, 6.3**

### Property 3: 浮点数与整数转换保持整数部分

*For any* 浮点数 f，将其转换为整数类型后，结果应等于 f 的整数部分（截断）。

**Validates: Requirements 1.4**

### Property 4: 大数类型与基础数值类型转换一致性

*For any* 基础数值类型值 v，将其转换为大数类型再转回基础类型，结果应等于原值（在精度范围内）。

**Validates: Requirements 1.5**

### Property 5: 复数转实数取实部

*For any* 复数 c，将其转换为实数类型后，结果应等于 c 的实部。

**Validates: Requirements 1.6**

### Property 6: 字符串与字节切片往返一致性

*For any* 有效的 UTF-8 字符串 s，ToString(ToBytes(s)) 应等于 s。

**Validates: Requirements 2.1**

### Property 7: 任意类型转字符串类成功

*For any* 支持的输入类型值 v（非 nil），ToString(v)、ToBytes(v)、ToStringer(v)、ToError(v) 都应成功（不返回错误）。

**Validates: Requirements 2.2, 2.3, 2.4, 2.5**

### Property 8: fmt.Stringer 和 error 输入使用字符串方法

*For any* 实现 fmt.Stringer 的类型值 v，将其转换为数值类型时，应使用 v.String() 的结果进行解析。

**Validates: Requirements 2.6, 2.7**

### Property 9: Duration 与数值类型往返一致性

*For any* time.Duration 值 d，ToDuration(ToInt64(d)) 应等于 d。

**Validates: Requirements 3.1**

### Property 10: 数值转布尔规则

*For any* 数值类型值 v，ToBool(v) 应在 v 非零时返回 true，v 为零时返回 false。

**Validates: Requirements 4.1, 4.4**

### Property 11: 布尔转数值规则

*For any* 布尔值 b，ToInt(b) 应在 b 为 true 时返回 1，b 为 false 时返回 0。

**Validates: Requirements 4.2**

### Property 12: 指针解引用正确性

*For any* 指针类型值 p 指向支持的类型，转换函数应正确解引用并转换底层值。

**Validates: Requirements 5.2**

### Property 13: nil 输入返回零值

*For any* 转换函数 ToXxx，ToXxx(nil) 应返回对应类型的零值且不返回错误。

**Validates: Requirements 6.1**

### Property 14: nil 大数指针返回错误

*For any* nil 的 *big.Int、*big.Float 或 *big.Rat 指针，转换函数应返回错误。

**Validates: Requirements 6.4**

### Property 15: 无效字符串返回错误

*For any* 无法解析为目标数值类型的字符串 s，ToXxxE(s) 应返回非 nil 错误。

**Validates: Requirements 6.5**

## 错误处理

### 错误类型

所有转换函数使用统一的错误格式：

```go
fmt.Errorf("unable to cast %#v of type %T to [目标类型]", a, a)
```

### 错误场景

| 场景 | 处理方式 |
|------|----------|
| 负数转无符号整数 | 返回错误 |
| nil 大数指针 | 返回错误 |
| 无效字符串格式 | 返回错误 |
| 不支持的类型 | 返回错误 |
| nil 输入 | 返回零值，无错误 |

## 测试策略

### 双重测试方法

本项目采用单元测试和属性测试相结合的方式：

- **单元测试**: 验证特定示例、边界情况和错误条件
- **属性测试**: 验证跨所有输入的通用属性

### 属性测试配置

- 使用 `github.com/leanovate/gopter` 作为属性测试库
- 每个属性测试至少运行 100 次迭代
- 每个测试用注释标记对应的设计属性
- 标记格式: **Feature: cast-extension, Property N: [属性描述]**

### 测试覆盖矩阵

测试应覆盖以下转换组合：

```
源类型 (25) × 目标类型 (22) = 550 种转换组合
```

重点测试：
1. 同类型转换（如 int → int64）
2. 跨类型转换（如 string → int）
3. 边界值（最大值、最小值、零值）
4. 错误情况（无效输入、溢出）

