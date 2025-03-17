# The input structure

OmniWOPE reads content from a JSON file. The JSON file consists of an array of posts. The posts will be published in the order they are listed in the file.

The file does not need to contain ALL posts ever published. Only those that were added or edited. OmniWOPE will retain mappings for all posts it ever published.

OmniWOPE will NEVER delete posts from the outputs.

```js
// wope.json
[
  {
    // URL is used as the unique post identifier
    "url": "https://myblog.com/canonical/post/url",
    "title": "Post title",
    "content": "Post content in Markdown",
    "date": "2025-03-15",
    "resources": [
      {
        // Path relative to "resources.base_path" in omniwope.yml
        "path": "path/to/file.jpg",
        "caption": "Markdown is supported",
        "type": "mime/type"
      }
    ]
    "tags": ["foo", "bar", "baz"]
  }
]
```
