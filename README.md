# OCCStructor - 中华人民共和国职业分类大典结构化工具

🏢 一个基于Go语言的职业分类数据解析工具，专门用于处理通过OCR技术从PDF转换而来的《中华人民共和国职业分类大典》Excel文件，并将其转换为结构化的JSON数据。

## ✨ 功能特性

- 📊 **智能Excel解析**: 自动识别并解析OCR转换的职业分类Excel文件
- 🌳 **分层树状结构**: 按大类→中类→小类→细类构建完整的职业分类层级体系
- 🤖 **AI智能处理**: 集成大模型，智能合并分割的职业名称，处理OCR错误
- 💾 **数据库存储**: 支持MySQL数据库，完整保存职业分类数据及层级关系
- 📤 **多格式导出**: 支持树状结构和扁平化两种JSON格式导出
- ⚙️ **灵活配置**: YAML配置文件 + 命令行参数，支持多环境部署
- 📋 **异常处理**: 自动记录解析异常，生成SQL模板便于人工修正
- 🔍 **正则优化**: 预编译正则表达式，提升解析性能

## 🚀 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- (可选) 通义千问API密钥

### 安装部署

```bash
# 1. 克隆项目
git clone https://github.com/solisamicus/occstructor.git
cd occstructor

# 2. 安装依赖
go mod tidy

# 3. 初始化数据库
mysql -u root -p < scripts/setup.sql

# 4. 配置环境变量(如果使用AI功能)
export DASHSCOPE_API_KEY="your_api_key_here"

# 5. 修改配置文件
cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml 设置数据库连接信息

# 6. 手动构建
go build -o bin/occstructor cmd/occstructor/main.go
go build -o bin/exportor cmd/exportor/main.go
```

### 基本使用

#### 1. 解析Excel并导入数据库

```bash
# 基本用法
./bin/occstructor -config configs/config.yaml

# 指定Excel文件
./bin/occstructor -excel "职业分类大典.xlsx"

# 直接运行(开发调试)
go run cmd/occstructor/main.go -config configs/config.yaml
```

#### 2. 导出JSON数据

```bash
# 导出树状JSON(默认)
./bin/exportor

# 导出扁平JSON
./bin/exportor -format=flat
# 完整参数示例
./bin/exportor \
  -format=tree \
  -output=data/occupations.json \
  -stats=true
```


## 📁 项目结构

```
OCCStructor/
├── cmd/            # 命令行工具
│ ├── occstructor/  # Excel解析导入工具
│ └── exportor/     # JSON导出工具
├── internal/       # 内部模块
│ ├── config/       # 配置管理
│ ├── model/        # 数据模型和树构建
│ ├── parser/       # Excel解析器核心
│ ├── repository/   # 数据访问层
│ └── service/      # 业务逻辑层
├── pkg/            # 公共包
│ └── database/     # 数据库连接
├── configs/        # 配置文件
├── scripts/        # 数据库脚本
├── logs/           # 日志文件
├── exports/        # 导出文件
└── docs/           # 文档
```
## ⚙️ 配置说明

### configs/config.yaml

```yaml
# 数据库配置
database:
  host: "localhost"
  port: 3306
  username: "root" 
  password: "your_password"
  database: "occupation_db"

# Excel文件配置
excel:
  filepath: "data/职业分类大典.xlsx"

# AI配置(可选)
ai:
  enabled: true
  api_key_env: "DASHSCOPE_API_KEY"
  base_url: "https://dashscope.aliyuncs.com/compatible-mode/v1/"
  model: "qwen-plus"
  temperature: 0.1

# 日志配置
logging:
  level: "info"
```

## 📊 数据格式说明

### 职业分类层级

| 层级 | 格式示例 | 说明 | GBM编码示例 |
|------|---------|------|------------|
| **大类** | `1` | 党的机关、国家机关... | `10` |
| **中类** | `1-01` | 中国共产党机关和基层组织负责人 | `10100` |
| **小类** | `1-01-00` | 中国共产党机关和基层组织负责人 | `10100` |
| **细类** | `1-01-00-01` | 中国共产党机关负责人 | - |

### JSON输出格式

#### 树状格式 (tree)
```json
{
  "data": [
    {
      "seq": "1",
      "gbm": "10", 
      "name": "党的机关、国家机关、群众团体和社会组织、企事业单位负责人",
      "level": 1,
      "children": [
        {
          "seq": "1-01",
          "gbm": "10100",
          "name": "中国共产党机关和基层组织负责人",
          "level": 2,
          "children": [
            // ... 更多子类别
          ]
        }
      ]
    }
  ],
  "stats": {
    "major_count": 8,
    "middle_count": 75, 
    "minor_count": 434,
    "detail_count": 1639
  },
  "exported_at": "2023-12-17T14:30:25Z",
  "total_records": 2156
}
```

#### 扁平格式 (flat)
```json
{
  "data": [
    {
      "seq": "1",
      "gbm": "10",
      "name": "党的机关、国家机关、群众团体和社会组织、企事业单位负责人", 
      "level": 1,
      "parent_seq": null
    },
    {
      "seq": "1-01", 
      "gbm": "10100",
      "name": "中国共产党机关和基层组织负责人",
      "level": 2,
      "parent_seq": "1"
    }
    // ... 更多记录
  ]
}
```

## 🔧 高级功能

### AI智能处理

当启用AI功能时，系统能够：
- 🔗 **智能合并**: 将分割的职业名称合并成完整标题
- ✂️ **智能分割**: 识别并分割组合的职业名称
- 🧹 **内容清理**: 自动移除OCR错误字符(如末尾的"L"、"S"等)
- 🔄 **格式标准化**: 统一职业名称格式

### 异常处理机制

当遇到代码与名称数量不匹配时：
1. **AI尝试修复** - 使用大模型智能处理
2. **记录详细日志** - 保存到 `logs/mismatch_*.log`
3. **生成SQL模板** - 便于手动修正数据

查看异常日志：
```bash
# 查看最新的不匹配日志
ls -la logs/mismatch_*.log
cat logs/mismatch_20250918_032621.log
```

### 性能优化

- ✅ **正则表达式预编译** - 程序启动时编译，提升解析速度
- ✅ **分层解析** - 按层级顺序处理，减少内存占用
- ✅ **批量数据库操作** - 事务批量插入，提升写入效率
- ✅ **智能备用处理** - AI不可用时自动降级到规则处理

## 📈 使用统计

目前已成功解析的职业分类数据规模：
- **8个大类** - 涵盖所有行业领域
- **79中类** - 细分行业分类
- **450小类** - 具体职业方向
- **1639细类** - 详细职业岗位