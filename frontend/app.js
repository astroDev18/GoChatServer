new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        email: null, // Email address used for grabbing an avatar
        username: null, // Our username
        joined: false // True if email and username have been filled in
    },
    // Now we take advantage of Vue's attribute of "created"
    // a function you define that handles anything you want to as soon as the Vue instance has been created
    //  
