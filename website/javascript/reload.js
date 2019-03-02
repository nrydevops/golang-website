window.onload = function() {

    if(performance.navigation.type == 2){

        if(window.location.pathname == "/login") {

            window.location.assign(window.location.hostname)
        }
        else {

            location.reload(true);
        }
    }
};

window.onunload = function(){};