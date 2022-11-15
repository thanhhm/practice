/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// $example on$
import org.apache.spark.graphx.GraphLoader
// $example off$
import org.apache.spark.sql.SparkSession

/**
 * A PageRank example on social network dataset
 * Run with
 * {{{
 * bin/run-example graphx.PageRankExample
 * }}}
 */
object PageRankExample {
  def main(args: Array[String]): Unit = {
    // Creates a SparkSession.
    val spark = SparkSession
      .builder
      .appName(s"${this.getClass.getSimpleName}")
      .getOrCreate()
    val sc = spark.sparkContext
    var productId = args.head.toDouble

    val t1 = System.nanoTime
    // $example on$
    // Load the edges as a graph
    // val graph = GraphLoader.edgeListFile(sc, "./data/graphx/followers.txt")
    val graph = GraphLoader.edgeListFile(sc, "./data/amazon/Amazon0601.txt")
    // Run PageRank
    val ranks = graph.pageRank(0.0001)
    val d1 = (System.nanoTime - t1) / 1e9d

    val t2 = System.nanoTime
    // Build connected component graph
    val cc = graph.triplets.collect {
      case t if t.srcId == productId => (t.dstId, t.srcId)
    }

    // Get the attributes of the top pagerank users
    val productWithRank = cc.join(ranks.vertices).map {
      case (id, (p, r)) => (id, r)
    }

    val d2 = (System.nanoTime - t2) / 1e9d

    // Print the result
    println("Build rerank graph exec time: ", d1)
    println("Filter neighborhood exec time:", d2)
    println(productWithRank.top(10)(Ordering.by(_._2)).mkString("\n"))
    // $example off$
    spark.stop()
  }
}
// scalastyle:on println
