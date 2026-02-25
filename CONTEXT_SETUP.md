# Context-Slim 使用说明

本项目已配置 context-slim 自动上下文优化工具。

## 自动更新机制

### Git Pre-commit Hook（已启用）

每次执行 `git commit` 时，会自动：
1. 扫描项目中的 .go 文件
2. 更新 `.context/` 目录下的索引和概览
3. 将更新的 .context/ 添加到本次提交

**Hook 位置：** `.git/hooks/pre-commit`

### 手动更新

如果需要手动更新上下文（不提交代码）：

```bash
# 方法 1：使用快捷脚本
./update-context.sh

# 方法 2：直接调用 context-slim
node /tmp/context-slim-temp/bin/cli.js update .
```

## 上下文文件说明

```
.context/
├── index.md              # L0: 项目文件索引（~15行）
├── CLAUDE.md             # 使用说明（已复制到根目录）
├── search-index.json     # 搜索索引
└── internal/
    └── downloader/
        └── _overview.md  # L1: 模块结构概览
```

- **L0（index.md）**：每个文件一行摘要，约 100 tokens
- **L1（_overview.md）**：每个模块的类型和函数签名，约 2k tokens
- **L2（源文件）**：完整源代码，按需加载

## 查询功能

搜索代码相关信息：

```bash
node /tmp/context-slim-temp/bin/cli.js query "downloader" .
node /tmp/context-slim-temp/bin/cli.js query "bilibili" .
```

## Watch 模式（可选）

实时监听文件变化并自动更新：

```bash
node /tmp/context-slim-temp/bin/cli.js watch .
```

## 注意事项

1. `.context/` 目录已添加到 `.gitignore`，不会提交到仓库
2. `CLAUDE.md` 会提交到仓库，供 AI 助手参考
3. 如果 context-slim 不可用，hook 会跳过更新（不影响提交）

## 性能提升

使用 context-slim 后：
- 传统方式：读取完整源码约 633 行（downloader 模块示例）
- L0 索引：仅 2 行描述
- L1 概览：仅 36 行签名
- **Token 节省 95%+**

## 卸载

如果不需要自动更新，删除 hook：

```bash
rm .git/hooks/pre-commit
```
