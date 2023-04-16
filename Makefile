
service-rewrite-pdf:
	go install github.com/mmcloughlin/podium@latest
	podium -output service-rewrite/rewrite.pdf http://127.0.0.1:3999/service-rewrite/rewrite.slide
