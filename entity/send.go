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
		a1, a2 := from.balance.amount, to.balance.amount
		return func() {
			if err != nil {
				from.balance.amount = a1
				to.balance.amount = a2
			}
		}
	}()

	err = from.balance.withdraw(amount)
	if err != nil {
		return err
	}
	to.balance.deposit(amount)
	return nil
}
