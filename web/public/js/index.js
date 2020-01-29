"use strict";
$(function() {
    var server = new WServer(
        location.hostname,
        3000,
        function(data) {
            const reader = new FileReader();
            reader.addEventListener('loadend', (e) => {
                var result = reader.result;

                var templateId = '#template-msg-row';
                var rowHtml = tpl(templateId, {"message" : result});

                $('#chat-container').append(rowHtml);
            });
            reader.readAsBinaryString(data);
        }
    );

    $('#send-frm').on('submit', function () {
        let el  = $('#send-msg');
        var value = el.val();
        el.val(null);
        console.log(value);
        server.send(value + "\n");
        return false;
    })

});

/**
 *
 * @param {string} host
 * @param {int} port
 * @param {function} receivedFunc
 * @constructor
 */
function WServer(host, port, receivedFunc) {
    let url = 'ws://' + host + ':' + port + '/ws';
    console.log(`Connect to "${url}"`);
    this.socket = new WebSocket(url);

    this.timerId = null;
    let self = this;

    this.socket.onopen = function(e) {
        console.log('[open] Connected.');
    };


    this.socket.onmessage = function(event) {
        console.log(`[message] Data received from server: ${event.data}`);
        receivedFunc(event.data)
    };

    this.socket.onclose = function(event) {
        if (event.wasClean) {
            console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
        } else {
            // e.g. server process killed or network down
            // event.code is usually 1006 in this case
            console.log('[close] Connection died');
        }
    };

    this.socket.onerror = function(error) {
        console.log(`[error] ${error.message}`);
    };

    /**
     * Send message to server.
     * @param {string} message
     */
    this.send = function (message) {
        this.socket.send(message);
    };
}

/**
 * @see https://golosay.net/simple-js-template-engine/
 *
 * @param {string} str
 * @param {object} data
 *
 * @returns {string}
 */
var tpl = function (str, data) {
    var name = [], value = [];
    var html = str.charAt(0) === '#' ? document.getElementById(str.substring(1)).innerHTML : str;
    if (typeof(data) === "object") {
        for (var k in data) {
            name.push(k);
            value.push(data[k]);
        }
    }
    var re = /{%([^%}]+)?%}/g, reExp = /(^( )?(var|if|for|else|switch|case|break|{|}|;))(?:(?=\()|(?= )|$)/g, code = 'var r=[];\n', cursor = 0, match;
    var add = function(line, js) {
        js? (code += line.match(reExp) ? line + '\n' : 'r.push(' + line + ');\n') :
            (code += line != '' ? 'r.push("' + line.replace(/"/g, '\\"') + '");\n' : '');
        return add;
    };
    while(match = re.exec(html)) {
        add(html.slice(cursor, match.index))(match[1], true);
        cursor = match.index + match[0].length;
    }
    add(html.substr(cursor, html.length - cursor));
    code += 'return r.join("");';
    return new Function(name, code.replace(/[\r\t\n]/g, '')).apply(this,value);
};


