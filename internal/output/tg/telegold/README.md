# Telegold - a Goldmark renderer for Telegram

If you want to publish regular Markdown into Telegram.

## The problem

Telegram allows [posting both in Markdown and in HTML](https://core.telegram.org/bots/api#formatting-options). But here's the catch: _neither is fully functional_.

If you choose to post in Markdown, you must escape all special characters not used for formatting. This is only feasible if you hand-craft the messages specifically for Telegram - or if you fully parse and re-render them, defeating the purpose of using Markdown instead of HTML.

If you choose to post in HTML, then only a strictly limited subset of tags is allowed. Worst of all, there is no line break formatting, so you can't reuse the same HTML for web publishing without adding line breaks.

## The solution

I wanted to be able to cross-post my blog posts into Telegram, using the same source. As I already use Markdown for posts, and Go for my telegram bot, I decided to extend the [Goldmark](https://github.com/yuin/goldmark) Markdown parser with a renderer that produces Telegram-compatible HTML.

## Usage

```go
md := goldmark.New(goldmark.WithRenderer(telegold.NewRenderer()))
md.Convert(source, &buf) // just use it as usual
```

As this is simply a renderer, you can combine it with other goldmark options.

---

(c) 2022 Leonid Shevtsov
