<?php
session_start();
require_once "../phpclient.php";
// Class for handling user responses
class Response {
    public function __construct()  {
    }
    public function setStatusCode(int $code) : Response {
        http_response_code($code);
        return $this;
    }

    public function toJson($data)  {
        header("Content-Type","application/json");
        echo json_encode($data);
    }
}
// class for handling user requests
class Request {
    private  $jsonInput;
    public function __construct() {
        // Gets json contents and convert to array
        $this->jsonInput = json_decode(file_get_contents('php://input'),true);
    }
    public function getJsonData($key)  {
        return $this->jsonInput[$key];
    }

    public function getRequest($key) {
        return $_GET[$key];
    }
}
// Function to login
function login() {
    $request = new Request();
    $response = new Response();
    $username = (string)$request->getJsonData("username");
    $password = (string)$request->getJsonData("password");
    $result = WebAuthorizationClient::Login($username, $password);
    if (!$result) {
        $response->setStatusCode(401)->toJson(['errorMessage' => 'Invalid username or password']);
    } else {
        $response->setStatusCode(200)->toJson(['success' => true]);
    }
}
// Function to refresh the token keeping users signed in
function refreshToken() {
    $response = new Response();
    if (WebAuthorizationClient::refreshToken()) {
        $response->setStatusCode(200)->toJson(['success' => true]);
    }
    else {
        $response->setStatusCode(401)->toJson(['success' => false]);
    }
}
// Function to get details of the logged in user
function checkme() {
    $response = new Response();
    if (!WebAuthorizationClient::checkAuth()) {
        $response->setStatusCode(401)->toJson(['success' => false]);
        return;
    }
    $result = WebAuthorizationClient::getUserDetails();
    $response->setStatusCode(200)->toJson($result);
}
// Function to logout
function logout() {
    $response = new Response();
    if (WebAuthorizationClient::logout()) {
        $response->setStatusCode(200)->toJson(['success' => true]);
    }
}

// Main entry point
function main() {
    $request = new Request();
    $route = $request->getRequest("page");
    switch ($route) {
        case 'login':
            login();
        break;
    
        case 'refreshToken':
            refreshToken();
        break;
    
        case 'checkMe': 
            checkme();
        break;
    
        case 'logout': 
            logout();
        break;
    }
}

main();