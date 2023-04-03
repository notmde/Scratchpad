const API_URL = "http://localhost:8081"
var userID
var pass

function myFunction() {
  var x = document.getElementById("password");
  if (x.type === "password") {
    x.type = "text";
  } else {
    x.type = "password";
  }
}

function dropDown() {
  var x = document.getElementsByClassName("dialog-box")[0].classList.toggle("show-box");
}

async function signup() {
  postData = {
    "_id": document.getElementById("text").value,
    "password": document.getElementById("password").value
  };

  try {
    const response = await fetch(API_URL + "/api/signup", {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(postData)
    });

    resp = await response.json()

    if (response.status == 400) {
      alert(resp["error"])
      return
    }

    if (response.status == 500) {
      console.log(resp["error"])
      return
    }

  } catch (error) {
    console.log('Failed to fetch data');
    return
  }

  var x = document.getElementById("actions-panel");
  x.children[0].remove()
  x.children[0].remove()

  var updateB = document.createElement("button");
  updateB.innerText = "Update";
  updateB.onclick = function () {
    updateData()
  };
  x.appendChild(updateB);

  var deleteB = document.createElement("button");
  deleteB.innerHTML = "Delete";
  deleteB.onclick = function () {
    deletion()
  };
  x.appendChild(deleteB);

  var logoutB = document.createElement("button");
  logoutB.innerHTML = "Logout";
  logoutB.onclick = function () {
    logout()
  };
  x.appendChild(logoutB);

  var name = document.getElementsByClassName("log")[0].children[0];
  userID = postData["_id"];
  name.innerText = postData["_id"];

  document.getElementById("text").readOnly = true;

  dropDown()
}

async function login() {
  postData = {
    "_id": document.getElementById("text").value,
    "password": document.getElementById("password").value
  };

  try {
    const response = await fetch(API_URL + "/api/login", {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(postData)
    });

    const resp = await response.json()

    clearCanvas()
    drawCreate(resp["canvas_data"])
    storeState.push(resp["canvas_data"])

    if (response.status == 400 || response.status == 401) {
      alert(resp["error"])
      return
    }

    if (response.status == 500) {
      console.log(resp["error"])
      return
    }

  } catch (error) {
    console.log('Failed to fetch data');
    return
  }

  var x = document.getElementById("actions-panel");
  x.children[0].remove()
  x.children[0].remove()

  var updateB = document.createElement("button");
  updateB.innerText = "Update";
  updateB.onclick = function () {
    updateData()
  };
  x.appendChild(updateB);

  var deleteB = document.createElement("button");
  deleteB.innerHTML = "Delete";
  deleteB.onclick = function () {
    deletion()
  };
  x.appendChild(deleteB);

  var logoutB = document.createElement("button");
  logoutB.innerHTML = "Logout";
  logoutB.onclick = function () {
    logout()
  };
  x.appendChild(logoutB);

  var name = document.getElementsByClassName("log")[0].children[0];
  userID = postData["_id"];
  name.innerText = postData["_id"];
  sessionStorage.setItem("password", postData["password"])

  document.getElementById("text").readOnly = true

  dropDown()
}

async function logout() {
  try {
    postData = {
      "canvas_data": storeState.pop()
    };

    const response = await fetch(API_URL + "/api/logout", {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(postData)
    });

    const resp = await response.json()

    if (resp.status == 500) {
      console.log(resp["error"])
      return
    }

  } catch (error) {
    console.log('Failed to fetch data');
    return
  }

  var x = document.getElementById("actions-panel");
  x.children[0].remove()
  x.children[0].remove()
  x.children[0].remove()

  var loginB = document.createElement("button");
  loginB.innerHTML = "Login";
  loginB.onclick = function () {
    login()
  };
  x.appendChild(loginB);

  var signupB = document.createElement("button");
  signupB.innerHTML = "Sign Up";
  signupB.onclick = function () {
    signup()
  };
  x.appendChild(signupB);

  var name = document.getElementsByClassName("log")[0].children[0];
  name.innerText = "Login";

  document.getElementById("text").readOnly = false;

  dropDown()
  clearCanvas()
}

async function updateData() {
  try {
    postData = {
      "password": document.getElementById("password").value
    };

    const response = await fetch(API_URL + "/api/update", {
      method: 'PATCH',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(postData)
    });

    const resp = await response.json()

    if (response.status == 500) {
      console.log(resp["error"])
      return
    }

  } catch (error) {
    console.log('Failed to fetch data');
    return
  }
}

async function deletion() {
  try {
    const response = await fetch(API_URL + "/api/delete", {
      method: 'DELETE',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      }
    });

    const resp = await response.json()

    if (response.status == 400) {
      alert(resp["error"])
      return
    }

  } catch (error) {
    console.log('Failed to fetch data');
    return
  }

  var x = document.getElementById("actions-panel");
  x.children[0].remove()
  x.children[0].remove()
  x.children[0].remove()

  var loginB = document.createElement("button");
  loginB.innerHTML = "Login";
  loginB.onclick = function () {
    login()
  };
  x.appendChild(loginB);

  var signupB = document.createElement("button");
  signupB.innerHTML = "Sign Up";
  signupB.onclick = function () {
    signup()
  };
  x.appendChild(signupB);

  var name = document.getElementsByClassName("log")[0].children[0];
  name.innerText = "Login";

  dropDown()
}

window.addEventListener("beforeunload", () => {
  sessionStorage.setItem("userID", userID);
  sessionStorage.setItem("canvasState", storeState.pop());
});

window.addEventListener("load", () => {
  if (sessionStorage.getItem("userID") != 'undefined' || sessionStorage.getItem("canvasState") != 'undefined') {
    var x = document.getElementById("actions-panel");
    x.children[0].remove()
    x.children[0].remove()

    var updateB = document.createElement("button");
    updateB.innerText = "Update";
    updateB.onclick = function () {
      updateData()
    };
    x.appendChild(updateB);

    var deleteB = document.createElement("button");
    deleteB.innerHTML = "Delete";
    deleteB.onclick = function () {
      deletion()
    };
    x.appendChild(deleteB);

    var logoutB = document.createElement("button");
    logoutB.innerHTML = "Logout";
    logoutB.onclick = function () {
      logout()
    };
    x.appendChild(logoutB);

    var name = document.getElementsByClassName("log")[0].children[0];
    name.innerText = sessionStorage.getItem("userID");

    document.getElementById("text").value = sessionStorage.getItem("userID")
    document.getElementById("text").readOnly = true
    document.getElementById("password").value = sessionStorage.getItem("password");
    
    clearCanvas()
  }

  if (sessionStorage.getItem("canvasState") != 'undefined') {
    drawCreate(sessionStorage.getItem("canvasState"))
  }
});