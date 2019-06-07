import { Component, OnInit } from "@angular/core";
import { RokService } from "../services/rok.service";
import { RokToken, EMPTY_TOKEN } from "src/app/uis/rok/utils/types";
import { Subscription } from "rxjs";
import { first } from "rxjs/operators";
import { Volume, Config } from "src/app/utils/types";
import {
  FormGroup,
  FormBuilder,
  Validators,
  FormControl
} from "@angular/forms";
import { NamespaceService } from "src/app/services/namespace.service";
import { KubernetesService } from "src/app/services/kubernetes.service";
import { Router } from "@angular/router";
import { getFormDefaults, initFormControls } from "src/app/utils/common";
import { createRokVolumeControl, addRokDataVolume } from "../utils/common";

@Component({
  selector: "app-rok-resource-form",
  templateUrl: "./rok-resource-form.component.html",
  styleUrls: ["./rok-resource-form.component.scss"]
})
export class RokResourceFormComponent implements OnInit {
  subscriptions: Subscription = new Subscription();
  token: RokToken = EMPTY_TOKEN;

  currNamespace = "";
  formCtrl: FormGroup;
  config: Config;

  ephemeral = false;
  defaultStorageclass = false;
  storageClasses: string[] = [];
  formReady = false;
  pvcs: Volume[] = [];

  constructor(
    private namespaceService: NamespaceService,
    private k8s: KubernetesService,
    private fb: FormBuilder,
    private router: Router,
    private rok: RokService
  ) {}

  ngOnInit() {
    // Init the form
    this.formCtrl = getFormDefaults();

    // Add labUrl control
    this.formCtrl.addControl("labUrl", new FormControl("", []));

    // Add the rok-url control
    const ws: FormGroup = this.formCtrl.get("workspace") as FormGroup;
    ws.controls.extraFields = this.fb.group({
      "rok-url": ["", [Validators.required]]
    });

    // Get the user's Rok Secret
    this.rok
      .getRokSecret("kubeflow")
      .pipe(first())
      .subscribe(token => {
        this.token = token;
      });

    // Update the form Values from the default ones
    this.k8s
      .getConfig()
      .pipe(first())
      .subscribe(config => {
        if (Object.keys(config).length === 0) {
          // Don't fire on empty config
          return;
        }

        this.config = config;
        this.initRokFormControls();
      });

    // Keep track of the selected namespace
    this.subscriptions.add(
      this.namespaceService.getSelectedNamespace().subscribe(namespace => {
        this.currNamespace = namespace;
        this.formCtrl.controls.namespace.setValue(this.currNamespace);

        this.updatePVCs(namespace);
      })
    );

    // Check if a default StorageClass is set
    this.k8s
      .getDefaultStorageClass()
      .pipe(first())
      .subscribe(defaultClass => {
        if (defaultClass.length === 0) {
          this.defaultStorageclass = false;
        } else {
          this.defaultStorageclass = true;
        }
      });

    // Get a list of existin StorageClasses
    this.k8s
      .getStorageClasses()
      .pipe(first())
      .subscribe(classes => {
        this.storageClasses = classes;
      });
  }

  ngOnDestroy() {
    // Unsubscriptions
    this.subscriptions.unsubscribe();
  }

  public updatePVCs(namespace: string) {
    this.subscriptions.add(
      this.k8s
        .getVolumes(namespace)
        .pipe(first())
        .subscribe(pvcs => {
          this.pvcs = pvcs;
        })
    );
  }

  public initRokFormControls() {
    // Sets the values from our internal dict. This is an initialization step
    // that should be only run once
    initFormControls(this.formCtrl, this.config);

    // Configure workspace control with rok-url
    this.formCtrl.controls.workspace = createRokVolumeControl(
      this.config.workspaceVolume.value
    );
    this.formCtrl
      .get("workspace")
      .get("path")
      .disable();

    // Add the data volumes
    const arr = this.fb.array([]);
    this.config.dataVolumes.value.forEach(vol => {
      // Create a new FormControl to append to the array
      addRokDataVolume(this.formCtrl, vol.value);
    });
  }

  public onSubmit() {
    this.formCtrl.updateValueAndValidity();
    const nb = JSON.parse(JSON.stringify(this.formCtrl.value));

    console.log(nb, this.formCtrl.valid);
    this.subscriptions.add(
      this.k8s
        .postResource(nb)
        .pipe(first())
        .subscribe(result => {
          if (result === "posted") {
            this.router.navigate(["/"]);
          } else if (result === "error") {
            this.updatePVCs(this.currNamespace);
          }
        })
    );
  }
}