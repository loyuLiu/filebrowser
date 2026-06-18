---
name: project-analysis
description: "Comprehensive analysis of a codebase: structure, tech stack, architecture, key files, and actionable summary"
---

# Project Analysis

Perform a thorough analysis of a project directory. Output a structured report covering overview, directory structure, tech stack, architecture, key files, and actionable observations.

## When to Use

- User asks to "analyze this project", "全面分析", "全面分析一下此項目", or similar
- Onboarding to an unfamiliar codebase
- Before making significant changes to understand impact

## Procedure

### Phase 1: Discovery (non-destructive exploration)

1. **List top-level structure**: `ls -la` the project root. Note: languages, frameworks, config files.
2. **Read entry points**: README.md, package.json, go.mod, requirements.txt, Makefile, Dockerfile, docker-compose.yml, compose.yaml — whatever exists.
3. **Map directory tree**: `find . -maxdepth 2 -type f | head -80` or similar to understand layout.
4. **Identify tech stack**: Infer from config files (package.json deps, go.mod, requirements.txt, etc.).

### Phase 2: Deep Dive (targeted reads)

5. **Read main source files**: The primary entry point, router/handler files, model definitions, config.
6. **Read tests**: Check for test structure and coverage patterns.
7. **Check CI/CD**: .github/workflows, Makefile targets, Taskfile.yml.
8. **Check database**: Migrations, schema files, ORM models.

### Phase 3: Report

Generate a structured report with these sections:

```markdown
## [Project Name] — 全面分析

### 1. 项目概览
| 项目 | 详情 |
|------|------|
| **名称** | ... |
| **技术栈** | ... |
| **定位** | 一句话描述项目用途 |

### 2. 目录结构
Tree-style overview of key directories and their purpose.

### 3. 架构概览
How the components connect (ASCII diagram if helpful).

### 4. 关键文件
List the most important files with brief descriptions.

### 5. 依赖与配置
External services, env vars, config files.

### 6. 代码质量观察
Patterns, potential issues, areas for improvement.

### 7. 总结
Actionable summary: what's working, what needs attention.
```

## Rules

- **Always in Chinese** (中文) for the report output, regardless of code language.
- **Non-destructive**: Only read files during analysis. Do not modify anything.
- **Structured output**: Use tables and headers for scanability.
- **Be specific**: Name actual files, functions, and line numbers where relevant.
- **Actionable**: End with concrete observations, not vague statements.
