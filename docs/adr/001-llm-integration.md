# ADR-001: LLM統合の設計方針

## ステータス
Accepted

## コンテキスト
ことわざの使い方を練習するロールプレイチャット機能にLLMを統合する必要がある。
複数のLLMプロバイダ（Anthropic Claude, OpenAI GPT等）が候補。

## 決定
- **Anthropic Claude Messages API** を採用
- LLMClient interfaceで抽象化し、プロバイダ交換可能な設計
- APIキーは環境変数 `LLM_API_KEY` で管理
- モデル指定は環境変数 `LLM_MODEL` で設定可能

## 理由
1. Claude APIはシステムプロンプトのサポートが強く、ことわざ教師のペルソナ設定に適する
2. Messages APIは構造化されており、会話履歴の管理が容易
3. interface抽象化により、テストでモック注入が可能
4. 将来的にプロバイダ切り替えやフォールバックが容易

## 影響
- LLM_API_KEY の設定が起動の必須条件
- LLMレスポンス時間がユーザー体験に直結（将来的にストリーミング対応を検討）
- テストではモックLLMClientを使用し、実際のAPI呼び出しは行わない
