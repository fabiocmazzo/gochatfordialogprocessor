var appVar = new Vue({
    el: '#app',
    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        email: null, // Email address used for grabbing an avatar
        fullname: null, // Our fullname
        username: null,
        uuid: null,
        joined: false, // True if email and username have been filled in
        questionTips: {}
    },
    computed: {
        showQuestionTips: function () {
            //return this.questionTips.length > 0;
            return true;
        }
    },


    created: function () {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function (e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
                + msg.username
                + '</div>'
                + self.videoParser(emojione.toImage(msg.message)) + '<br/>'; // Parse emojis
            self.chatContent += self.parseForm(msg.form);
            self.questionTips = msg.questionTipList;
            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        });
    },
    methods: {
        sendMessage: function (message) {
            this.ws.send(
                JSON.stringify({
                        email: this.email,
                        username: this.username,
                        hidden: false,
                        clientAppId: this.uuid,
                        message: $('<p>').html(message).text() // Strip out html
                    }
                ));
        },
        videoParser: function(html){
            var pattern1 = /(?:http?s?:\/\/)?(?:www\.)?(?:vimeo\.com)\/?(.+)/g;
            var pattern2 = /(?:http?s?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=)?(.+)/g;
            var pattern3 = /([-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?(?:jpg|jpeg|gif|png))/gi;
            if(pattern1.test(html)){
                var replacement = '<iframe width="420" height="345" src="//player.vimeo.com/video/$1" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>';

                var html = html.replace(pattern1, replacement);
            }
            if(pattern2.test(html)){
                var replacement = '<iframe width="420" height="345" src="http://www.youtube.com/embed/$1" frameborder="0" allowfullscreen></iframe>';
                var html = html.replace(pattern2, replacement);
            }
            if(pattern3.test(html)){
                var replacement = '<a href="$1" target="_blank"><img class="sml" src="$1" /></a><br />';
                var html = html.replace(pattern3, replacement);
            }
            return html;
        },

        send: function () {
            if (this.newMsg != '') {
                this.sendMessage(this.newMsg);
                this.newMsg = ''; // Reset newMsg
            }
        },

        join: function () {
            var validated = true;
            if (!this.email) {
                 Materialize.toast('Por favor, informe seu email', 2000);
                 validated = false;
             }
             if (!this.fullname) {
                 Materialize.toast('Por favor, informe seu nome completo', 2000);
                 validated = false;
             }

             if(!validated) {
                return;
             }



            // Enviar um PING para o processador
            // para avisar que o usuário entrou
            // Inclui o additionalInformation para informar
            // os dados antes de qualquer processmento
            this.uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
                var r = Math.random()*16|0, v = c == 'x' ? r : (r&0x3|0x8);
                return v.toString(16);
            });

            var additionalInfo =  [{name: "email", datatype:"EMAIL", value: this.email}, {name: "fullname", datatype:"FULLNAME", value: this.fullname}];
            this.ws.send(
                JSON.stringify({
                        email: this.email,
                        username: this.username,
                        additionalInformationList : additionalInfo,
                        hidden: true,
                        clientAppId: this.uuid,
                        message: "PING_PING_PING_PING_PING"
                    }
                ));



            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html("Você").text();
            this.joined = true;
        },

        gravatarURL: function (email) {

            if (!email) {
                emailHash = "emailnapreenchido";
            } else {
                emailHash = CryptoJS.MD5(email);
            }
            return 'http://www.gravatar.com/avatar/' + emailHash;
        },

        parseForm: function (formExpression) {
            if (!formExpression) {
                return "";
            }

            // Começar os testes por exemplo se é Boolean Sim:Não
            formHtml  =  "<a @click=\"sendMessage('Sim')\" class='waves-effect waves-light btn'><i class='material-icons left'>thumb_up</i>Sim</a>";
            formHtml  += "<a @click=\"sendMessage('Não')\" class='waves-effect waves-light btn'><i class='material-icons left'>thumb_down</i>Não</a>";;

            return formHtml;
        }


    }
});