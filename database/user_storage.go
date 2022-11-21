package database

var temporaryUserStorage = make([]*UserStorage, 0)

type UserStorage struct {
	userId string
	data   map[string]interface{}
}

// AddUserStorage adds a new user storage to the temporary storage
func AddUserStorage(userId string, data map[string]interface{}) {
	temporaryUserStorage = append(temporaryUserStorage, &UserStorage{
		userId: userId,
		data:   data,
	})
}

func GetUserById(userId string) *UserStorage {
	for _, user := range temporaryUserStorage {
		if user.userId == userId {
			return user
		}
	}
	return nil
}

func (u *UserStorage) Get(key string) interface{} {
	return u.data[key]
}

func (u *UserStorage) Set(key string, value interface{}) {
	u.data[key] = value
}

func (u *UserStorage) Delete(key string) {
	delete(u.data, key)
}

func (u *UserStorage) Clear() {
	u.data = make(map[string]interface{})
}

func (u *UserStorage) GetUserId() string {
	return u.userId
}

func (u *UserStorage) SetUserId(userId string) {
	u.userId = userId
}

func (u *UserStorage) GetData() map[string]interface{} {
	return u.data
}

func (u *UserStorage) SetData(data map[string]interface{}) {
	u.data = data
}
