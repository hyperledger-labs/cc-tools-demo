'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const { readFile } = require('fs').promises;

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.assets = []
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        this.data = roundArguments.filePath ? (await readFile(roundArguments.filePath)).toString('base64') : roundArguments.data

        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
    }

    async submitTransaction() {
        const randomId = this.generateRandomString(8)
        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'createAsset',
            invokerIdentity: 'User1',
            contractArguments: ['{"asset":[{"@assetType": "benchmarkAsset", "data": "'+this.data+'", "id": "'+randomId+'"}]}'],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(myArgs);
        this.assets.push(randomId)
    }

    async cleanupWorkloadModule() {
        for (const assetID of this.assets) {
            console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'deleteAsset',
                invokerIdentity: 'User1',
                contractArguments: ['{"key": {"@assetType":"benchmarkAsset","id":"'+assetID+'"}}'],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
        }
    }

    generateRandomString(length) {
        const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        let result = "";
        for (let i = 0; i < length; i++) {
            const randomIndex = Math.floor(Math.random() * charset.length);
            result += charset.charAt(randomIndex);
        }
        return result;        
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;