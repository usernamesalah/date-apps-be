package request

type UserRegister struct {
	Name        string `json:"name" valid:"required"`
	Email       string `json:"email" valid:"email,required"`
	Password    string `json:"password" valid:"required"`
	PhoneNumber string `json:"phone_number" valid:"numeric,optional"`
}

type UserLogin struct {
	Email       string `json:"email" valid:"email,required"`
	Password    string `json:"password" valid:"required"`
	PhoneNumber string `json:"phone_number" valid:"numeric,optional"`
}
