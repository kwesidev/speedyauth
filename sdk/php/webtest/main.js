let refreshTokenInterval = null;
const form = document.getElementById('upload-form');
const fileInput = document.getElementById('file-input');
const uploadButton = document.getElementById('upload-button');
/*
uploadButton.addEventListener('click', (event) => {
  event.preventDefault();
  const file = fileInput.files[0];
  const chunkSize = 1024 * 1024; // 1 MB
  let start = 0;
  while (start < file.size) {
    const chunk = file.slice(start, start + chunkSize);

    const xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:3000/users/chunkUpload');
    xhr.setRequestHeader('Content-Type', 'application/octet-stream');
    xhr.setRequestHeader('X-File-Name', encodeURIComponent(file.name));
    xhr.setRequestHeader('X-File-Size', file.size);
    xhr.setRequestHeader('X-File-Start', start);
    let token = window.localStorage.getItem('token');
    if (token) {
        xhr.setRequestHeader("token", token);
    }
    console.log(chunk);
    xhr.send(chunk);
    start += chunkSize;
  }
});
*/
// Get cookie 
function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function login() {
    let username, password;
    username = document.getElementById('username').value;
    password = document.getElementById('password').value;
    apiRequest({
        path: 'backend.php?page=login',
        method: 'POST',
        data: { username: username, password: password },
        onSuccess: function (response) {
            if (response.success) {
                refreshTokenInterval = setInterval(refreshToken, 1200000);
                init();
            }
        },
        onFailure: function (resultCode, resultText) {
            let result = JSON.parse(resultText);
            document.getElementById('errorMessage').innerHTML = result.error || result.errorMessage;
        }
    });
}
function logout() {
    apiRequest({
        path: 'backend.php?page=logout',
        method: 'POST',
        onSuccess: function (response) {
            if (response.success) {
                window.localStorage.clear();
                window.location.reload();
                clearInterval(refreshTokenInterval);
            }
        },
        onFailure: function(response){
            window.localStorage.clear();
            window.location.reload();
            clearInterval(refreshTokenInterval);
        }
    });
}

function refreshToken() {
    apiRequest({
        path: 'backend.php?page=refreshToken',
        method: 'GET',
        onSuccess: function (response) {
        },
        onFailure: function(response) {
            logout();
        }
    });
}
function apiRequest(request) {
    let bodyRequest, xhr, method, path, csrfToken;
    xhr = new XMLHttpRequest();
    if (request.hasOwnProperty('method')) {
        method = request.method;
    }
    else {
        method = 'GET';
    }
    if (request.hasOwnProperty('path')) {
        path = request.path;

    }
    if (path == null || '') {
        throw Error('Path is required');
    }
    xhr.open(method, path, true);
    csrfToken = getCookie('csrfToken');
    if (csrfToken) {
        xhr.setRequestHeader("csrfToken", csrfToken);
    }
    //Send the proper header information along with the request
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = () => { // Call a function when the state changes.
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            let response = JSON.parse(xhr.responseText);
            if (request.hasOwnProperty('onSuccess')) {
                request.onSuccess(response);
            }
        }
        else if (xhr.readyState == XMLHttpRequest.DONE && xhr.status !== 200) {
            if (request.hasOwnProperty('onFailure') && request.onFailure !== null) {
                request.onFailure(xhr.status, xhr.responseText);
            }
        }
    }
    // xhr.setRequestHeader("Content-type", "application/json");
    if (request.hasOwnProperty('data')) {
        bodyRequest = JSON.stringify(request.data);
    }
    xhr.send(bodyRequest);
}

function init(){
    apiRequest({
        path: 'backend.php?page=checkMe',
        method: 'GET',
        onSuccess: function (response) {
            document.getElementById('loginPage').style.display = 'none';
            document.getElementById('welcomePage').style.display = 'block';
            document.getElementById('welcomeMessage').innerHTML = '<h1>' + response.username + '</h1>';
        }
    });
}