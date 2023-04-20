<?php
session_start();
if (isset($_SERVER['HTTP_ORIGIN'])) {
    header("Access-Control-Allow-Origin: ".$_SERVER['HTTP_ORIGIN']);
}else {
    header("Access-Control-Allow-Origin: *");
}
header("Access-Control-Allow-Methods: GET,HEAD,OPTIONS,POST,PUT");
header('Access-Control-Allow-Credentials: true');
header("Access-Control-Allow-Headers: content-type,token,csrfToken");
header("Content-Type: application/json; charset=UTF-8");
require_once "../WebAuthServerClient.php";

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
// Csrtoken checks
class CsrfToken {
    public static function init() {
        if (!isset($_COOKIE['csrfToken'])) {
            CsrfToken::generateToken();
        }
    }
    public static  function generateToken() {
        $randomString = bin2hex(random_bytes(12));
        setcookie('csrfToken',$randomString,0,"/","",isset($_SERVER['HTTPS']),false);
    }
    /**
     * Function to check csrftokens
     * @return bool
     */
    public static function check() : bool {
        $csrfToken = (string)$_SERVER['HTTP_CSRFTOKEN'];
        $response = new Response();
        if (isset($csrfToken) && isset($_COOKIE['csrfToken'])) {
            if ($csrfToken == $_COOKIE['csrfToken']) {
                return true;
            }
            else {
                $response->setStatusCode(400)->toJson(['errorMessage' => 'Cross site fogery detected']);
                die();
            }
        }
        $response->setStatusCode(400)->toJson(['errorMessage' => 'CSR Token needed']);
        die();
    }
}
// Determine which Request method to use for specific functions
function enforceRequestMethod(string $type)  {
    if ($_SERVER['REQUEST_METHOD'] !== $type)  {
        $response = new Response();
        $response->setStatusCode(200)->toJson(['message' => 'This method not allowed']);
        die();
    }
}
// Function to login
function login() {
    enforceRequestMethod('POST');
    CsrfToken::check();
    $request = new Request();
    $response = new Response();
    $username = (string)$request->getJsonData("username");
    $password = (string)$request->getJsonData("password");
    $result = WebAuthServerClient::login($username, $password);
    if (!$result) {
        $response->setStatusCode(401)->toJson(['errorMessage' => 'Invalid username or password']);
    } else {
        $response->setStatusCode(200)->toJson(['success' => true]);
    }
}
// Function to refresh the token keeping users signed in
function refreshToken() {
    enforceRequestMethod('GET');
    CsrfToken::check();
    $response = new Response();
    if (WebAuthServerClient::refreshToken()) {
        $response->setStatusCode(200)->toJson(['success' => true]);
    }
    else {
        $response->setStatusCode(401)->toJson(['success' => false]);
    }
}
// Function to get details of the logged in user
function checkme() {
    enforceRequestMethod('GET');
    CsrfToken::check();
    $response = new Response();
    if (!WebAuthServerClient::checkAuth()) {
        $response->setStatusCode(401)->toJson(['success' => false]);
        return;
    }
    $result = WebAuthServerClient::getUserDetails();
    $response->setStatusCode(200)->toJson($result);
}
// Function to logout
function logout() {
    enforceRequestMethod('POST');
    CsrfToken::check();
    $response = new Response();
    if (WebAuthServerClient::logout()) {
        $response->setStatusCode(200)->toJson(['success' => true]);
    }
}
// Main entry point
function main() {
    CsrfToken::init();
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
            // Logout and generate new csrftoken not to keep token for too long
            logout();
            CsrfToken::generateToken();
        break;
    }
}

main();