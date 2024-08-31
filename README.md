Pock is a small language to perform simple calculations and comparisons. I named
it "Pock" because it reminded me of the small expressions I would type in my
dad's [Pocket computer](https://en.wikipedia.org/wiki/Pocket_computer) back in
the 80s.

![Casio FX-850P pocket computer](https://c1.staticflickr.com/9/8191/8082061808_bffa4c5220_b.jpg)

# Getting Started

To run Pock, type:

```shell
go install github.com/loderunner/pocklang/cmd/pock
pock
```

## Loading state

You can load a JSON file as an immutable state for the interpreter, and refer to
the values of the state in your pock expressions.

```json
{
  "hello": "world",
  "THX": 1138,
  "foo": {
    "bar": {
      "baz": true
    }
  }
}
```

```
❯ pock --state state.json
Pock v0.0.0
> hello
"world"
> THX
1138
> foo.bar.baz
true
```
