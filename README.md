# This example is to be used with the [Browserless V1 template](https://railway.app/template/browserless)

Create a reference variable on your Railway service that you deploy your app to

```shell
BROWSER_WS_ENDPOINT=${{Browserless.BROWSER_WS_ENDPOINT}}
```

</br>

Then use `os.Getenv("BROWSER_WS_ENDPOINT")` in code.

### Before

```go
// create context
ctx, cancel := chromedp.NewContext(
    context.Background(),
)

defer cancel()
```

### After

```go
// create allocator context
allocatorContext, cancel := chromedp.NewRemoteAllocator(
    context.Background(),
    os.Getenv("BROWSER_WS_ENDPOINT"),
    chromedp.NoModifyURL,
)

defer cancel()

// create context
ctx, cancel := chromedp.NewContext(allocatorContext)

defer cancel()
```

The rest of your Go code remains the same with no other changes required.