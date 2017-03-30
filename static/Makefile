all: virtualjoystick.js materialize

# materialize
materialize.js:
	curl https://cdnjs.cloudflare.com/ajax/libs/materialize/0.98.0/js/materialize.min.js > materialize.js
materialize.css:
	curl https://cdnjs.cloudflare.com/ajax/libs/materialize/0.98.0/css/materialize.min.css > materialize.css
materialize: materialize.js materialize.css

# virtualjoystick.js
virtualjoystick.js:
	curl https://raw.githubusercontent.com/jeromeetienne/virtualjoystick.js/master/virtualjoystick.js > virtualjoystick.js

clean:
	rm materialize.css materialize.js virtualjoystick.js

test: all
