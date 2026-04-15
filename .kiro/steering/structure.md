# Project Structure

## Organization Philosophy

仕様駆動開発に基づくプロジェクト構成。仕様書とステアリングを先に整備し、実装を後から追加するアプローチ。

## Directory Patterns

### Steering
**Location**: `.kiro/steering/`
**Purpose**: プロジェクト全体のルールとコンテキストをAIに提供する
**Example**: `product.md`, `tech.md`, `structure.md`

### Specifications
**Location**: `.kiro/specs/`
**Purpose**: 個別機能の仕様書（要件定義・設計・タスク）を管理する
**Example**: `.kiro/specs/{feature-name}/requirements.md`

### Settings
**Location**: `.kiro/settings/`
**Purpose**: テンプレートとルール定義（メタデータ）

## Naming Conventions

- **ドキュメント**: kebab-case（例: `feature-name`）
- **仕様ディレクトリ**: 機能名をkebab-caseで命名

## Code Organization Principles

- ソースコードの構成は設計フェーズで決定する
- 仕様書に基づいた実装を行う

---
_Document patterns, not file trees. New files following patterns shouldn't require updates_
