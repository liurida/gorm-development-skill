# Key Concepts for Belongs To

A `belongs to` association sets up a one-to-one connection with another model. Each instance of the declaring model "belongs to" one instance of the other model.

## Foreign Key

To define a `belongs to` relationship, a foreign key must exist. The default foreign key uses the owner's type name plus its primary key field name.

For the example above, to define the `User` model that belongs to the `CreditCard`, the foreign key should be `CreditCardID`.

## Overriding Foreign Key

GORM allows you to override the foreign key with the `foreignKey` tag.

```go
type User struct {
  gorm.Model
  Name         string
  MyCreditCardID uint
  CreditCard   CreditCard `gorm:"foreignKey:MyCreditCardID"`
}
```
