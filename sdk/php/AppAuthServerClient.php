<?php 
/**
 * This library communicates directly to the Authentication Server 
 * This is useful when making a single page app or a mobile and want to store tokens on the app 
 */
class AppAuthServerClient {
    const AuthServerUrl = 'http://localhost:8080';
    /**
     * Function to login
     * @param string username
     * @param string password
     */
    public static function login(string $username, string $password) : array  {
        $curl = curl_init();
        $jsonData = json_encode([
            'username' => $username,
            'password' => $password,
        ]);
        curl_setopt($curl, CURLOPT_URL, AppAuthServerClient::AuthServerUrl."/api/auth/login");
        curl_setopt($curl, CURLOPT_POST, 1);
        curl_setopt($curl, CURLOPT_HTTPHEADER, array('Content-Type:application/json'));
        curl_setopt($curl, CURLOPT_POSTFIELDS, $jsonData);
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
        $result = json_decode(curl_exec($curl), true);
        if (array_key_exists("token", $result)) {
            $result['success'] = true;
            return $result;
        }
        return $result;
    }
    
    /**
     * Function to refresh the token
     * @param {String} $refreshToken
     * @return token
     */
    public static function refreshToken(string $refreshToken): array {
        $curl = curl_init();
        $jsonData = json_encode([
            'refreshToken' => $refreshToken
        ]);
        curl_setopt($curl, CURLOPT_URL, AppAuthServerClient::AuthServerUrl . "/api/auth/tokenRefresh");
        curl_setopt($curl, CURLOPT_HTTPHEADER, array('Content-Type:application/json'));
        curl_setopt($curl, CURLOPT_POST, 1);
        curl_setopt($curl, CURLOPT_POSTFIELDS, $jsonData);
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
        $result = json_decode(curl_exec($curl), true);
        if (array_key_exists("token", $result)) {
            $result['success'] = true;
            return $result;
        }
        return $result;
    }
    /**
     * Function to check the login status
     * @return bool
     */
    public static function checkAuth(): bool {
        $token = (string)$_SERVER['HTTP_TOKEN'];
        if (!isset($token)) {
            return false;
        }
        $curl = curl_init();
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($curl, CURLOPT_URL, AppAuthServerClient::AuthServerUrl . "/api/user");
        curl_setopt($curl, CURLOPT_HTTPHEADER, array('token:'.$token));
        $result = json_decode(curl_exec($curl),true);
        if (array_key_exists("id",$result)) {
            return true;
        }
        else {
            return false;
        }
    }
    /**
     * Function to logout
     * @param {String} $refreshToken
     */
    public static function logout(string $refreshToken) : bool {
        $curl = curl_init();
        $jsonData = json_encode([
            'refreshToken' => $refreshToken
        ]);
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($curl, CURLOPT_URL, AppAuthServerClient::AuthServerUrl . "/api/auth/logout");
        curl_setopt($curl, CURLOPT_POST, 1);
        curl_setopt($curl, CURLOPT_POSTFIELDS, $jsonData);
        $result = json_decode(curl_exec($curl), true);
        if (array_key_exists("success", $result)) {
                return $result['success'];
        }
        else {
            return false;
        }
    }
    /**
     * Function to get user details
     * @return {Object}
     */
    public static function getUserDetails() : array {
        $token = (string)$_SERVER['HTTP_TOKEN'];
        if (!isset($token)) {
            return [];
        }
        $curl = curl_init();
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($curl, CURLOPT_URL, AppAuthServerClient::AuthServerUrl . "/api/user");
        curl_setopt($curl, CURLOPT_HTTPHEADER,array('token:'.$token));
        $result = json_decode(curl_exec($curl), true);
        if (array_key_exists("id",$result)) {
            return $result;
        }
        else {
            return [];
        }
    }

    /**
     * Function to check if user is admin
     * @param {String} $token 
     * @return bool
     */
    public static function isUserAdmin(string $token) : bool {
        $curl = curl_init();
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($curl, CURLOPT_URL, AppAuthServerClient::AuthServerUrl . "/api/user");
        curl_setopt($curl, CURLOPT_HTTPHEADER,array('token:'.$token));
        $result = json_decode(curl_exec($curl),true);
        if (array_key_exists("id",$result)) {
            if (in_array("ADMIN",$result['roles'])) {
                return true;
            }
            else {
                return false;
            }
        }
        else {
            return false;
        }
    }
}