package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job"`
}

// Declarando minha array de dados a serem manipulados
var people []Person

func main() {
	router := mux.NewRouter() // Cria um router usando a lib Mux
	port := ":1337"           // Define a porta em que vai rodar a minha aplicaço

	// Populando minha array de dados
	people = append(people, Person{ID: "1", Name: "Jhon", Age: 20, Job: "Professor"})
	people = append(people, Person{ID: "2", Name: "Maria", Age: 23, Job: "Scientist"})
	people = append(people, Person{ID: "3", Name: "Mark", Age: 25, Job: "Developer"})
	people = append(people, Person{ID: "4", Name: "Clark", Age: 33, Job: "Attorney"})
	people = append(people, Person{ID: "5", Name: "Lilian", Age: 35, Job: "Dentist"})

	router.HandleFunc("/people", getPeople).Methods("GET")            // Rota para receber todos os dados da minha array
	router.HandleFunc("/people", createPerson).Methods("POST")        // Rota para criar um novo registro na minha array
	router.HandleFunc("/people/{id}", getPerson).Methods("GET")       // Rota para receber um unico dado da minha array a partir do ID
	router.HandleFunc("/people/{id}", updatePerson).Methods("PUT")    // Rota para atualizar um dado da minha array a partir do ID
	router.HandleFunc("/people/{id}", deletePerson).Methods("DELETE") // Rota para excluir um dado da minha array a partir do ID

	fmt.Println("Starting server at port", port) // Informa no console que o servidor esta rodando
	log.Fatal(http.ListenAndServe(port, router)) // Inicia o servidor usando a lib http e passando o Router criado na linha 24
}

// Todas as funçoes Rest devem receber dois parametros:
// writer -> que eh o metodo para retornarmos algo a partir da requisiçao
// request -> que eh o conteudo que recebemos no momento da requisiçao
func updatePerson(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json") // Iniciamos todas as requisiçoes informando no Header que o Content-Type eh application/json

	params := mux.Vars(request) // Aqui usei a lib do Mux para receber o request da requisiçao

	for index, item := range people { // Aqui iteramos nossos dados para encontrar o elemento com o ID correspondente ao ID do request
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...) // Assim que encontramos o ID correspondente removemos o mesmo da array
			var person Person                                    // Criamos um novo dado que sera inserido na array no lugar do element removido
			_ = json.NewDecoder(request.Body).Decode(&person)    // Passamos os dados do Body do request para dentro da variavel criada
			person.ID = params["id"]                             // Colocamos o ID correto no novo dado criado
			people = append(people, person)                      // Adicionamos o novo elemento na array de dados
			json.NewEncoder(writer).Encode(person)               // Usamos o writer para devolver ao cliente uma resposta que contem o novo dado cadastrado
			return
		}
	}
}

func deletePerson(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, item := range people {
		if item.ID == params["id"] { // Encontra o ID passado por parametro
			people = append(people[:index], people[index+1:]...) // Remove o dado com o ID correspondente da array de dados
			break
		}
	}
	json.NewEncoder(writer).Encode(people) // Retorna para o cliente o dado que foi removido
}

func getPerson(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for _, item := range people {
		if item.ID == params["id"] { // Encontra o ID passado por parametro
			json.NewEncoder(writer).Encode(item) // Devolve para o cliente o dado correspondente ao ID
			return
		}
	}
}

func createPerson(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var person Person                                 // Criar variavel que sera nosso novo dado
	_ = json.NewDecoder(request.Body).Decode(&person) // Pega o Body do request e decodifica para a variavel criada anteriormente
	person.ID = strconv.Itoa(rand.Intn(1000000))      // Gera um numero aleatorio para o ID do novo dado criado
	people = append(people, person)                   // Adiciona o novo dado na array de dados
	json.NewEncoder(writer).Encode(person)            // Devolve para o cliente uma resposta contendo o novo dado criado
}

func getPeople(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(people) // Retorna toda a array de dados para o cliente
}
