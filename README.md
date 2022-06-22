# Emoji 🤣

📖 ⏹ 🔙 🛬 ❓ 🏅 🛥 🕟 🗾 🔻 🆚 ↔️ 🍒 🎭 🎌 🍈 🛥 💞 👹 🌛 🌺 🕚 🐓 3️⃣ 🐍 ♓️ 🗓

## Installation

```shell
go get github.com/sari3l/emoji
```

## Usage

```go
emoji.Activity()        // Enable
emoji.InActivity()      // Disable
```

## Example

```go
import "github.com/sari3l/emoji"

func main() {
    emoji.Activity()
    fmt.Println("Enable: ::beer: Hello Emoji!!!:alien:")
    emoji.InActivity()
    fmt.Println("Disable: ::beer: Hello Emoji!!!:alien:")
}
```

Output:

```shell
╰─○ go build -gcflags="all=-l" main.go && ./main
Enable: :🍺 Hello Emoji!!!👽
Disable: ::beer: Hello Emoji!!!:alien:
```

## Notes

1. Sometimes fails to patch a function if inlining is enabled. Try running your tests with inlining disabled, for example: `-gcflags="all=-l"`. The same command line argument can also be used for build.
