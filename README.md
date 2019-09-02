# curriculum-mapper

*** DEPRECATED (RETIRED) ***

*This repository is obsolete and retired (archived). This is an unmantained repository. In particular, note that it is not being updated with security patches, including those for dependent libraries.*

A machine learning tool to map between curricula

NOTE: This is experimental and proof-of-concept code

This tool uses keyword-based document classification to align one curriculum to another curriculum. It identifies distinctive keywords in each source curriculum item text, and uses those as a document classifier. The tool then runs the text of each target curriculum item past the classifier, and extracts scores for how well the target item overlaps in keywords with each source curriculum item. 

Binary distributions of the code are available in the build/ directory.

For somewhat more on the approach taken, see https://github.com/nsip/curriculum-mapper/wiki/Design-Approach

The document classification approach is used in this instance to align curricula, but the documnt classification approach applies to any text; it can just as well be applied to resource descriptions, or to lesson plans.

To test the approach, we provide the Year 7 and 8 Science curriculum text from the Australian Curriculum as the source curriculum, and the Stage 4 (Years 7–8) NSW Science Syllabus as the target curriculum. Note the following idiosyncracies in the data sets:

* The Australian Curriculum includes content descriptions, and elaborations of content descriptions (which are examples). The document classifier is trained on both sets of text.
* The NSW Syllabus differentiates outcomes (which are competencies) and content. The two are in principle autonomous; but outcomes and content are aligned, by belonging to the same Strand and Substrand: each Strand/Substrand combination has a single content text, and one or two outcomes. So the content text and the outcomes text can be aligned.
  * The Skills strand (which is competency-oriented) has one outcome per Strand + Substrand.
  * The Knowledge and Understanding strand (which is content-oriented) has one or two outcomes per Strand + Substrand; so the context text has broader coverage than each outcome.
  * The Values and Attitudes strand has no associated content text.
* The tool attempts to align 17 NSW Syllabus outcomes to 37 Australian Curriculum content descriptions
* Being competencies, we expect the NSW outcome text to align more poorly to the Australian Curriculum content descriptions than the associated NSW content text does.
* The NSW content text partly overlaps with the Australian Curriculum, and nominates aligned Australian Curriculum content descriptions.

We use Australian Curriculum content descriptions mention in the NSW content text as a control. The document classifier processes the outcome text and the associated content text separately, trying to align each to the Australian Curriculum.

* If the best Australian Curriculum content descriptions match for the NSW outcome text and for the associated NSW content text both match one of the  Australian Curriculum content descriptions nominated in the content text, then there is an alignment score of 2.
* If the best Australian Curriculum content descriptions match for either the NSW outcome text or for the associated NSW content text matches one of the Australian Curriculum content descriptions nominated in the content text, then there is an alignment score of 1.
* If the best Australian Curriculum content descriptions match for neither the NSW outcome text nor for the associated NSW content text matches one of the Australian Curriculum content descriptions nominated in the content text, then there is an alignment score of 0.

The results of the alignment are shown in the Excel document included in the distribution, which presents the output of the document classifiers as generated by this code. 

* Columns are Australian Curriculum content descriptions 
* Rows are NSW Syllabus outcomes
* The score in each cell is the alignment score for the outcome text, followed by the alignment score for the associated content text
* "%" flags an Australian Curriculum content description mentioned in the NSW content text; we expect the document classifier to identify these successfully. In the Excel file out.xslx, these are in red text.
* "#" flags the best match to a Australian Curriculum content description in the NSW outcome text and the NSW content text. In the the Excel file out.xslx, these are cells with red borders.

The results are as follows:

* The -VA outcomes (Values and Attitudes) have no associated content text, so the accuracy of the match cannot be confirmed.
* There are 8 Knowledge and Understanding outcomes (the first 8). 5 have an alignment score of 2, 3 have an alignment score of 1. (So average alignment score: 1.625.) In all three outcomes with an alignment score of 1, it was the outcome text rather than the content text that led to a misalignment, as predicted.
* There are 6 -WS (Skills) outcomes. 3 have an alignemnt score of 2, 2 have an alignment score of 1, 1 has an alignment score of 0. The average alignment score is 1.333: this is poorer than the alignment score for content-based outcomes, as predicted. In one of the three instances (SC4_9WS), the misalignment was in the content text match rather than the outcome text match; but the difference of the winning content text match from the nearest correct match was small (-2091 vs -2096).

These results indicate that keyword-based alignment performs better than chance, and can be the basis of a useful recommender.
