# JWT Authentication in Golang (Gorilla Mux) 🔐
A secure JWT-based authentication system built in Golang using Gorilla Mux and GORM (PostgreSQL).

✅ User Registration & Login with hashed passwords (bcrypt).
✅ JWT Token Generation & Verification.
✅ Middleware for Protected Routes using JWT.
✅ PostgreSQL Database Integration with GORM.



# 📌 Installation & Setup
1️⃣ Clone the Repository
git clone https://github.com/sangeeth200494/JWT-AUTH_Golang.git

cd JWT-AUTH_Golang


2️⃣ Install Dependencies
go mod tidy


3️⃣ Configure Environment Variables
Create a .env file and add the following:
DB_HOST=localhost
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=jwt_auth_db
DB_PORT=5432
SECRET_KEY=your_secret_key


4️⃣ Run the Application
go run main.go



# 🛠 API Endpoints
Method	Endpoint	Description	Authorization
POST	/register	   Register a new user	        ❌ No Token
POST	/login	     Authenticate user & get JWT 	❌ No Token
GET	  /protected	 Access protected content	    ✅ Requires Token



# 🔑 JWT Authentication Flow
1️⃣ User Registers/Login → Receives a JWT Token.
2️⃣ Client Stores the Token (e.g., Local Storage, HTTP Headers).
3️⃣ Client Sends Requests with Authorization: Bearer <token>.
4️⃣ Server Verifies the Token before allowing access.



# 🛠 Technologies Used
Golang (Gorilla Mux for routing)
GORM (PostgreSQL ORM)
JWT (JSON Web Tokens)
bcrypt (Password hashing)


# 📜 License
This project is open-source and licensed under the MIT License.


# 📞 Contact
# 👤 Sangeeth Jayaraj
# 📧 Email: sangeethvv554@gmail.com
# 🔗 GitHub: github.com/sangeeth200494

