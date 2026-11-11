package nosuppress

//nolint:errcheck // want `suppression comments are not honored`
func a() {}

// nosec G104 // want `suppression comments are not honored`
func b() {}

//lint:ignore U1000 unused // want `suppression comments are not honored`
func c() {}

//edil:ignore nomagic because // want `suppression comments are not honored`
func d() {}

// A normal comment about linting is fine.
func e() {}
