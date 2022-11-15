export SPH=/home/thanhhm/Downloads/tool/spark/spark-3.3.1-bin-hadoop3

simple-app:
	${SPH}/bin/spark-submit \
	--class "SimpleApp" \
	--master local[2] \
	target/scala-2.13/simple-project_2.13-1.0.jar

submit-pagerank:
	${SPH}/bin/spark-submit \
	--class "PageRankExample" \
	--master local[2] \
	target/scala-2.12/simple-project_2.12-1.0.jar 190629

submit-connected:
	${SPH}/bin/spark-submit \
	--class "ConnectedComponentsExample" \
	--master local[2] \
	target/scala-2.12/simple-project_2.12-1.0.jar

submit-compre:
	${SPH}/bin/spark-submit \
	--class "ComprehensiveExample" \
	--master local[2] \
	target/scala-2.12/simple-project_2.12-1.0.jar 183331
