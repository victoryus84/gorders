allprojects {
    repositories {
        google()
        mavenCentral()
    }
}

// Correct Kotlin DSL syntax to force local project path
rootProject.layout.buildDirectory.set(rootProject.file("../build"))

subprojects {
    val newBuildDir = rootProject.layout.buildDirectory.dir(project.name).get().asFile
    project.layout.buildDirectory.set(newBuildDir)
}

subprojects {
    project.evaluationDependsOn(":app")
}

tasks.register<Delete>("clean") {
    delete(rootProject.layout.buildDirectory)
}