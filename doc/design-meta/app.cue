package flyb

source: "baldrick-seer.design"
name: "baldrick-seer"
modules: ["design"]

graphIntegrityPolicy: {
  missingNode:              "error"
  orphanNode:               "warning"
  duplicateNoteName:        "error"
  unknownRelationshipLabel: "ignore"
  crossReportReference:     "allow"
}
