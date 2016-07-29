# choる

~~ラーメン次郎~~[mattn/cho](https://github.com/mattn/cho)インスパイア系

## Example

```go
c := choru.New()
if i, v := c.Choose([]string{"foo", "bar", "hoge", "fuga"}); i >= 0 {
	fmt.Printf("%d => %s\n", i, v)
}
```

## Screen Shot
![capture](capture.gif)

## Keybind

|      key     | operation |
|:------------:|:---------:|
|      `j`     |    down   |
|      `k`     |     up    |
|      `g`     |    top    |
|      `G`     |   bottom  |
| `q`, `<ESC>` |    quit   |
