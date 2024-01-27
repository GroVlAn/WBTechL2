package pattern

import "fmt"

/*
Преймущества:
	- Избавляет от множества больших условных операторов машины состояний.
	- Концентрирует в одном месте код, связанный с определённым состоянием.
	- Упрощает код контекста.
Недостатки:
	- Может неоправданно усложнить код, если состояний мало и они редко меняются.
*/

// UserContext - контекст пользователя, с текущим состоянием
type UserContext struct {
	State UserState
}

// RegularUserState - состояние пользователя не имеющего прав
type RegularUserState struct{}

// AdminUserState - состояние администратора с дополнительными правами
type AdminUserState struct{}

func (rus *RegularUserState) GrantPermission(permission string) error {
	fmt.Printf("permission %s is grant\n", permission)
	return nil
}

func (rus *RegularUserState) RevokePermission(permission string) error {
	fmt.Printf("permission %s is revoke\n", permission)
	return nil
}

func (aus *AdminUserState) GrantPermission(permission string) error {
	fmt.Printf("permission %s is grant\n", permission)
	return nil
}

func (aus *AdminUserState) RevokePermission(permission string) error {
	fmt.Printf("permission %s is revoke\n", permission)
	return nil
}

// ChangeState - метод меняющий состояние контекста
func (uc *UserContext) ChangeState(state UserState) {
	uc.State = state
}

// UserState - интерфейс состояния пользователя
type UserState interface {
	GrantPermission(permission string) error
	RevokePermission(permission string) error
}
