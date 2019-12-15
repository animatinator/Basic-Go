package data

type Database struct {
	Data map[string]string
}

func NewDb() Database {
	return Database {make(map[string]string)}
}

func (d *Database) Read(key string) <-chan string {
	c := make(chan string, 1);
	
	go func() {
		defer close(c)

		// TODO: Delays.
		// TODO: Random errors.
		v, ok := d.Data[key]
		// TODO: Do something if a key is missing.
		if ok {
			c <- v
		}
	}()
	
	return c;
}

func (d *Database) Write(key, value string) <- chan struct{} {
	c := make(chan struct{});
	
	go func() {
		defer close(c)
		
		d.Data[key] = value
	}()
	
	return c;
}