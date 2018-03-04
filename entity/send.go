package entity

// send locks balances of from and to and move amount
// send is called from Pay or AcceptInvoice
func send(from, to *User, amount Amount) error {
	from.balance.mu.Lock()
	defer from.balance.mu.Unlock()

	to.balance.mu.Lock()
	defer to.balance.mu.Unlock()

	// rollback users if err was occurred
	var err error
	defer func() func() {
		b1, b2 := *from.balance, *to.balance
		return func() {
			if err != nil {
				from.balance = &b1
				to.balance = &b2
			}
		}
	}()

	err = from.balance.withdraw(amount)
	if err != nil {
		return err
	}
	err = to.balance.deposit(amount)
	if err != nil {
		return err
	}
	return nil
}
