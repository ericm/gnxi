import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { environment } from '../environment';
import { FileService } from '../file.service';
import { Targets } from '../models/Target';
import { TargetService } from '../target.service';

@Component({
  selector: 'app-devices',
  templateUrl: './devices.component.html',
  styleUrls: ['./devices.component.css']
})
export class DevicesComponent implements OnInit {
  targetForm: FormGroup;
  selectedTarget: any;
  targetList: Targets;

  constructor(private http: HttpClient, private targetService: TargetService, private formBuilder: FormBuilder, private fileService: FileService, private snackbar: MatSnackBar) { }

  ngOnInit(): void {
    this.targetForm = this.formBuilder.group({
      targetName: ['', Validators.required],
      address: ['', Validators.required],
      ca: ['', Validators.required],
      cakey: ['', Validators.required],
    });
    this.selectedTarget = {};
    this.getTargets();
  }

  getTargets(): void {
    this.targetService.getTargets().subscribe(targets => {
      if (targets != null) {
        this.targetList = targets;
      } else {
        this.targetList = {};
      }
    });
  }

  setTarget(targetForm): void {
    this.http.post(`${environment.apiUrl}/target/${targetForm.targetName}`, targetForm).subscribe(
      (res) => {
        this.targetList[targetForm.targetName] = {
          address: targetForm.address,
          ca: targetForm.ca,
          cakey: targetForm.cakey,
        };
        this.selectedTarget = this.targetList[targetForm.targetName];
        this.snackbar.open("Saved", "", {duration: 2000});
    },
      (error) => console.error(error),
    );
  }

  deleteTarget(): void {
    let name = this.targetForm.get("targetName").value;
    this.fileService.deleteFile(this.selectedTarget.ca).subscribe((res) => {
      console.log(res);
    }, error => console.error(error));
    this.fileService.deleteFile(this.selectedTarget.cakey).subscribe((res) => {
      console.log(res);
    }, error => console.error(error));
    this.targetService.delete(name).subscribe(res => {
      this.selectedTarget = {};
      this.targetForm.reset();
      delete this.targetList[name];
      this.snackbar.open("Deleted", "", {duration: 2000});
    }, error => console.error(error))
  }

  setSelectedTarget(targetName: string): void {
    this.selectedTarget = this.targetList[targetName];
    if (this.selectedTarget === undefined) {
      this.selectedTarget = {};
      this.targetForm.reset();
      return;
    }
    this.targetForm.setValue({
      targetName,
      address: this.selectedTarget.address,
      ca: this.selectedTarget.ca,
      cakey: this.selectedTarget.cakey,
    });
  }

  addCa(caFileName: string): void {
    this.targetForm.patchValue({
       ca: caFileName,
    });
  }

  addCaKey(keyFileName: string): void {
    this.targetForm.patchValue({
      cakey: keyFileName,
    });
  }

  get targetName(): AbstractControl {
    return this.targetForm.get('targetName');
  }
  get targetAddress(): AbstractControl {
    return this.targetForm.get('address');
  }
}
