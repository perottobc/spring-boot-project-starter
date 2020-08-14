package upgrade

import (
	"co-pilot/pkg/springio"
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func SpringBoot(model *pom.Model) error {
	springRootInfo, err := springio.GetRoot()
	if err != nil {
		return err
	}

	currentVersion, err := getSpringBootVersion(model)
	if err != nil {
		return err
	}

	latestVersion, err := ParseVersion(springRootInfo.BootVersion.Default)
	if err != nil {
		return err
	}

	if currentVersion.IsDifferentFrom(latestVersion) {
		msg := fmt.Sprintf("outdated spring-boot version [%s => %s]", currentVersion.ToString(), latestVersion.ToString())
		if IsMajorUpgrade(currentVersion, latestVersion) {
			log.Warnf("major %s", msg)
		} else if !latestVersion.IsReleaseVersion() {
			log.Warnf("%s | not release", msg)
		} else {
			log.Info(msg)
		}

		return updateSpringBootVersion(model, latestVersion)
	} else {
		log.Infof("Spring boot is the latest version [%s]", latestVersion.ToString())
	}

	return nil
}

func getSpringBootVersion(model *pom.Model) (JavaVersion, error) {
	// check parent
	if model.Parent.ArtifactId == "spring-boot-starter-parent" {
		return ParseVersion(model.Parent.Version)
	}

	// check dependencyManagement
	if model.DependencyManagement != nil {
		dep, err := model.DependencyManagement.Dependencies.FindArtifact("spring-boot-dependencies")
		if err != nil {
			return JavaVersion{}, nil
		}
		version, err := model.GetDependencyVersion(dep)
		if err != nil {
			return JavaVersion{}, nil
		}
		return ParseVersion(version)
	}

	return JavaVersion{}, errors.New("could not extract spring boot version information")
}

func updateSpringBootVersion(model *pom.Model, newVersion JavaVersion) error {
	// check parent
	if model.Parent.ArtifactId == "spring-boot-starter-parent" {
		model.Parent.Version = newVersion.ToString()
		return nil
	}

	// check dependencyManagement
	if model.DependencyManagement != nil {
		dep, err := model.DependencyManagement.Dependencies.FindArtifact("spring-boot-dependencies")
		if err != nil {
			return err
		} else {
			return model.SetDependencyVersion(dep, newVersion.ToString())
		}
	}

	return errors.New("could not update spring boot version to " + newVersion.ToString())
}
