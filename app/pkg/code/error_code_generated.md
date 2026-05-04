# 错误码

！！系统错误码列表，由 `codegen -type=int -doc` 命令生成，不要对此文件做任何更改。

## 功能说明

如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：

```json
{
  "code": 100101,
  "message": "Database error"
}
```

上述返回中 `code` 表示错误码，`message` 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。

## 错误码列表

系统支持的错误码列表如下：

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
| ErrRecordNotFound | 101101 | 404 | Record not found |
| ErrMessageQuery | 101102 | 500 | Failed to query Message from Database |
| ErrMessageCreate | 101103 | 500 | Failed to create Message in Database |
| ErrGoodsNotFound | 100501 | 404 | Goods not found |
| ErrCategoryNotFound | 100502 | 404 | Category not found |
| ErrEsUnmarshal | 100503 | 500 | Elasticsearch unmarshal error |
| ErrBannerNotFound | 100504 | 404 | Banner not found |
| ErrBrandNotFound | 100505 | 404 | Brand not found |
| ErrCategoryBrandNotFound | 100506 | 404 | CategoryBrand not found |
| ErrGoodsImageNotFound | 100507 | 404 | GoodsImage not found |
| ErrJsonUnmarshal | 100508 | 500 | JSON unmarshal error |
| ErrInventoryNotFound | 100601 | 404 | Inventory not found |
| ErrInvSellDetailNotFound | 100602 | 404 | Inventory sell detail not found |
| ErrInvNotEnough | 100603 | 400 | Inventory not enough |
| ErrOptimisticRetry | 100604 | 500 | Optimistic lock retry limit exceeded |
| ErrShopCartItemNotFound | 100701 | 404 | ShopCart item not found |
| ErrSubmitOrder | 100702 | 500 | Failed to submit order |
| ErrNoGoodsSelect | 100703 | 400 | No goods selected |
| ErrOrderNotFound | 100704 | 404 | Order not found |
| ErrOrderStatus | 100705 | 500 | Failed to update order status |
| ErrInvalidParameter | 100706 | 400 | Invalid request parameter |
| ErrUserNotFound | 100401 | 404 | User not found |
| ErrUserPasswordIncorrect | 100402 | 401 | User password incorrect |
| ErrCodeNotExist | 100403 | 404 | Verification code not exist |
| ErrCodeInCorrect | 100404 | 400 | Verification code incorrect |
| ErrUserAlreadyExists | 100405 | 400 | User already exists |
| ErrSmsSend | 100406 | 500 | Failed to send SMS |
| ErrForbidden | 100407 | 403 | User privilege insufficient |

