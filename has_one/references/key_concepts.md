# Key Concepts for Has One

A `has one` association is also a one-to-one connection, but it sets up the relationship in the opposite direction of a `belongs to`.

## Foreign Key

For a `has one` relationship, the owned entity must contain the foreign key. The default foreign key is the owner's type name plus its primary key.

In the example above, the `CreditCard` model has a `UserID` field, which is the foreign key.

## Overriding Foreign Key

You can override the foreign key with the `foreignKey` tag.

```go
type User struct {
  gorm.Model
  CreditCard CreditCard `gorm:"foreignKey:MyUserID"`
}

type CreditCard struct {
  gorm.Model
  Number   string
  MyUserID uint
}
```
