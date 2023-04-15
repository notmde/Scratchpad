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

    //userID = postData["_id"]
    // isLoggedin = true
    // sessionStorage.setItem("isLoggedin", isLoggedin);
    // document.querySelector("#login").innerText = "Logout"
    // document.querySelector("#login").onclick = logout

    // dropDown()

    login()
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

    //userID = postData["_id"]
    isLoggedin = true
    sessionStorage.setItem("isLoggedin", isLoggedin);
    document.querySelector("#login").innerText = "Logout"
    document.querySelector("#login").onclick = logout

    dropDown()
}

async function logout() {
    try {
        postData = {
            "_id": document.getElementById("text").value,
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

    isLoggedin = false
    sessionStorage.setItem("isLoggedin", isLoggedin)
    document.querySelector("#login").innerText = "Login"
    document.querySelector("#login").onclick = dropDown

    clearCanvas()
}