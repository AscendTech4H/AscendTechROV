all: virtualjoystick.js jquery.js materialize

# materialize
materialize.js:
	curl https://cdnjs.cloudflare.com/ajax/libs/materialize/0.98.0/js/materialize.min.js > materialize.js
materialize.css:
	curl https://cdnjs.cloudflare.com/ajax/libs/materialize/0.98.0/css/materialize.min.css > materialize.css
materialize: materialize.js materialize.css

# jquery
jquery.js:
	curl https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js > jquery.js

# virtualjoystick.js
virtualjoystick.js:
	curl https://raw.githubusercontent.com/jeromeetienne/virtualjoystick.js/master/virtualjoystick.js > virtualjoystick.js

clean:
	rm jquery.js materialize.css materialize.js virtualjoystick.js

