package maven

import (
	"errors"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

func Merge(from *pom.Model, to *pom.Model) error {
	if err := mergeDependencies(from, to); err != nil {
		log.Errorln(err)
	}

	if err := mergeManagementDependencies(from, to); err != nil {
		log.Errorln(err)
	}

	if err := mergePlugins(from, to); err != nil {
		log.Errorln(err)
	}

	return nil
}

func mergeDependencies(from *pom.Model, to *pom.Model) error {
	if from.Dependencies == nil || to.Dependencies == nil {
		return errors.New("dependencies is nil")
	}

	for _, fromDep := range from.Dependencies.Dependency {
		var hasDependency = false
		for _, toDep := range to.Dependencies.Dependency {
			if fromDep.GroupId == toDep.GroupId && fromDep.ArtifactId == toDep.ArtifactId {
				hasDependency = true
			}
		}
		if !hasDependency {
			log.Infof("inserting dependency %s:%s into project", fromDep.GroupId, fromDep.ArtifactId)
			to.Dependencies.Dependency = append(to.Dependencies.Dependency, fromDep)
		}
	}

	return nil
}

func mergeManagementDependencies(from *pom.Model, to *pom.Model) error {
	if from.DependencyManagement == nil || to.DependencyManagement == nil {
		return errors.New("dependencyManagement is nil")
	}

	for _, fromDepMan := range from.DependencyManagement.Dependencies.Dependency {
		var hasManagementDependency = false
		for _, toDepMan := range to.DependencyManagement.Dependencies.Dependency {
			if fromDepMan.GroupId == toDepMan.GroupId && fromDepMan.ArtifactId == toDepMan.ArtifactId {
				hasManagementDependency = true
			}
		}
		if !hasManagementDependency {
			log.Infof("inserting management dependency %s:%s into project", fromDepMan.GroupId, fromDepMan.ArtifactId)
			to.DependencyManagement.Dependencies.Dependency = append(to.DependencyManagement.Dependencies.Dependency, fromDepMan)
		}
	}

	return nil
}

func mergePlugins(from *pom.Model, to *pom.Model) error {
	if from.DependencyManagement == nil || to.DependencyManagement == nil {
		return errors.New("build.plugin is nil")
	}

	for _, fromPlugin := range from.Build.Plugins.Plugin {
		var hasPlugin = false
		for _, toPlugin := range to.Build.Plugins.Plugin {
			if fromPlugin.GroupId == toPlugin.GroupId && fromPlugin.ArtifactId == toPlugin.ArtifactId {
				hasPlugin = true
			}
		}
		if !hasPlugin {
			log.Infof("inserting plugin %s:%s into project", fromPlugin.GroupId, fromPlugin.ArtifactId)
			to.Build.Plugins.Plugin = append(to.Build.Plugins.Plugin, fromPlugin)
		}
	}

	return nil
}