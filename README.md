 weather-app-final


 Weather App


This project is a backend application designed for managing and serving weather data. The application follows the principles of Clean Architecture, ensuring modularity, scalability, and testability. It interacts with the OpenWeatherMap API to fetch weather data and stores it in a MongoDB database for efficient retrieval and updates.


 Features


- Fetch current weather data from OpenWeatherMap API.
- Store and retrieve weather data from a MongoDB database.
- Update weather data via API calls.
- RESTful endpoints for weather-related operations.
- Graceful error handling and structured logging.


 Technologies Used


- Programming Language: Go (Golang)
- Database: MongoDB
- External API: OpenWeatherMap
- Environment Management: dotenv
- HTTP Server: Built-in Go `net/http`
- Dependency Management: Go Modules


 Project Structure


weather-app-final/
├── cmd/                 Entry point for the application
├── config/              Configuration and environment management
├── internal/
│   ├── handler/         HTTP handlers (controllers)
│   ├── models/          Data models
│   ├── repository/      Database operations
│   ├── service/         Business logic
├── .env                 Environment variables
├── go.mod               Dependency management
├── go.sum               Dependency checksums
└── main.go              Main application entry point


 Installation


1. Clone the repository:
  git clone https://github.com/mashinolol/weather-app-final.git
  cd weather-app


2. Install dependencies:
  go mod tidy


3. Create a `.env` file in the root directory and configure the following variables:
  MONGO_URI=mongodb:
  DATABASE_NAME=weatherdb
  API_KEY= yourapikey
  BASE_URL=https://api.openweathermap.org/data/2.5/weather


4. Run the application:
  go run main.go


 Endpoints


 1. Get Weather
Endpoint: `/weather?city={city_name}` 
Method: `GET` 
Description: Retrieves weather data for the specified city.


Response:
{
 "city": "TestCity",
 "temp": 25.0,
 "description": "Sunny"
}


 2. Update Weather
Endpoint: `/weather?city={city_name}` 
Method: `PUT` 
Description: Updates the weather data for the specified city.


Response:
Weather data updated successfully



